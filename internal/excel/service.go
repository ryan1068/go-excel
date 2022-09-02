package excel

import (
	"context"
	"cst/internal/pkg/config"
	"cst/pkg/request"
	"cst/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var workerNum = 4 //同时运行的工作协程

var pageSize int64 = 100 //获取数据的每页数据量

type Service struct {
	cfg     *config.Config
	redis   *redis.Client
	mongodb *mongo.Client
	ctx     context.Context
}

type Result struct {
	row []string
	res []gjson.Result
	err error
}

func (s *Service) createTask(form *ImportForm) (string, error) {

	f := func() (string, [][]string, string, error) {
		filePath, err := s.downloadExcel(form.FilePath, form.TaskType)
		if err != nil {
			return "", nil, "", err
		}

		f, err := excelize.OpenFile(filePath)
		if err != nil {
			return "", nil, "", err
		}
		// Get all the rows in the Sheet1.
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			return "", nil, "", err
		}

		rows = s.handleRows(rows)
		mapping := strings.Split(form.Mapping, ",")
		if len(mapping) != len(rows[0]) {
			return "", nil, "", errors.New("上传Excel中表头格式设置不正确")
		}

		apiHost := s.getApiHost(form.ApiHost)
		if apiHost == nil {
			return "", nil, "", errors.New("apiHost传参不正确")
		}

		taskId := utils.RandStringBytes(10)
		return filePath, rows, taskId, nil
	}

	filePath, rows, taskId, logId, err := s.log(f, form)
	if err != nil {
		return "", err
	}
	go s.importExcel(form, filePath, rows, taskId, logId)

	return taskId, nil
}

// 导入excel文件
func (s *Service) importExcel(form *ImportForm, filePath string, rows [][]string, taskId string, logId interface{}) error {
	newFile := excelize.NewFile()
	index := newFile.NewSheet("Sheet1")
	if err := s.generateHeader(rows[0], newFile); err != nil {
		return err
	}

	totalRow := len(rows) - 1
	cacheKey := s.cacheKey(form.ApiParams["group_id"], form.ApiParams["area_id"], form.TaskType)
	s.redis.Set(s.ctx, cacheKey, taskId, time.Second*3600*2)
	s.redis.HSet(s.ctx, taskId, "count", totalRow)
	s.redis.Expire(s.ctx, taskId, time.Second*3600*2)

	taskChan := make(chan []string, totalRow)
	resChan := make(chan Result)
	doneChan := make(chan struct{}, totalRow)

	go s.worker(form, taskChan, resChan, doneChan)
	go s.producer(rows, taskChan)

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	failedStartRow := 1 //导入失败日志开始写入Excel的行数
	for {
		select {
		case ch, ok := <-resChan:
			if !ok {
				if failedStartRow > 1 {
					// 如果有失败记录，生成失败记录Excel
					errorFile, err := s.createErrorFile(filePath)
					if err != nil {
						return err
					}
					basePath := strings.TrimLeft(errorFile, "/www/runtime/go-excel/")
					s.redis.HSet(s.ctx, taskId, "filePath", fmt.Sprintf("%v:%v/%v", s.cfg.Intranet.Ip, s.cfg.Application.Port, basePath))
					colStr, _ := utils.ConvertNumToCol(len(rows[0]) + 1)
					newFile.SetColWidth("Sheet1", "A", colStr, 20)
					if err := newFile.SaveAs(errorFile); err != nil {
						return err
					}
				}
				s.redis.HSet(s.ctx, taskId, "isDone", 1)
				s.updateLog(logId,
					bson.D{{"$set", bson.D{
						{"progress", s.getProgress(taskId)},
						{"status", 1},
					}}})
				return nil
			}

			if ch.err != nil {
				// 写入失败记录到Excel
				failedStartRow++
				for col, value := range ch.row {
					axis := s.getAxis(col+1, failedStartRow)
					newFile.SetCellValue("Sheet1", axis, value)
				}
				axis := s.getAxis(len(rows[0])+1, failedStartRow)
				newFile.SetCellValue("Sheet1", axis, ch.err)
				newFile.SetActiveSheet(index)
				s.redis.HIncrBy(s.ctx, taskId, "failCount", 1)
			} else {
				s.redis.HIncrBy(s.ctx, taskId, "successCount", 1)
			}
		case <-ticker.C:
			if len(doneChan) == totalRow {
				close(resChan)
			}
			s.updateLog(logId,
				bson.D{{"$set", bson.D{
					{"progress", s.getProgress(taskId)},
				}}})
		}
	}
	return nil
}

// 生产任务
func (s *Service) producer(rows [][]string, taskChan chan<- []string) {
	for i, row := range rows {
		if i == 0 {
			continue
		}
		taskChan <- row
	}
	close(taskChan)
}

// 消费任务
func (s *Service) worker(form *ImportForm, taskChan <-chan []string, resChan chan<- Result, doneChan chan struct{}) {
	for i := 0; i < workerNum; i++ {
		go func() {
			for {
				row, ok := <-taskChan
				if !ok {
					return
				}
				dataMap := s.buildApiDataMap(form, row)
				res, err := s.requestApi(form, dataMap)
				resChan <- Result{
					row: row,
					res: res,
					err: err,
				}
				doneChan <- struct{}{}
			}
		}()
	}
}

// 获取数据所在excel表格的坐标
func (s *Service) getAxis(col int, row int) string {
	colStr, _ := utils.ConvertNumToCol(col)
	axis := colStr + strconv.Itoa(row)
	return axis
}

// 获取需要处理的数据行
func (s *Service) handleRows(rows [][]string) [][]string {
	var validRows [][]string
	for _, row := range rows {
		if row == nil {
			continue
		}
		//if len(row) != len(rows[0]) {
		//	continue
		//}
		validRows = append(validRows, row)
	}
	return validRows
}

// 生成错误日志文件
func (s *Service) createErrorFile(filePath string) (string, error) {
	errlogDir := "/www/runtime/go-excel/static/" + time.Now().Format("20060102") + "/errlog"
	if err := os.MkdirAll(errlogDir, 0755); err != nil {
		return "", err
	}

	basePath := path.Base(filePath)
	return errlogDir + "/" + basePath, nil
}

// 生成excel表头
func (s *Service) generateHeader(row []string, newFile *excelize.File) error {
	for col, value := range row {
		axis := s.getAxis(col+1, 1)
		newFile.SetCellValue("Sheet1", axis, value)
	}

	axis := s.getAxis(len(row)+1, 1)
	newFile.SetCellValue("Sheet1", axis, "失败原因")
	return nil
}

// 下载oss文件到本地
func (s *Service) downloadExcel(path, taskType string) (string, error) {
	url := s.cfg.Oss.Url + "/" + path

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	filePath, err := s.createDownloadFile(url, taskType)
	if err != nil {
		return "", err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// 生成下载文件
func (s *Service) createDownloadFile(url, taskType string) (string, error) {
	uploadDir := "/www/runtime/go-excel/static/" + time.Now().Format("20060102") + "/upload"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	base := path.Base(url)
	ext := path.Ext(base)
	filePath := fmt.Sprintf("%v/%v-%v%v", uploadDir, taskType, time.Now().Unix(), utils.RandStringBytes(5)+ext)
	return filePath, nil
}

// 获取导入任务缓存key
func (s *Service) cacheKey(groupId, areaId, taskType string) string {
	if groupId != "" {
		return fmt.Sprintf("group:%v:type:%v", groupId, taskType)
	}
	if areaId != "" {
		return fmt.Sprintf("area:%v:type:%v", areaId, taskType)
	}
	return taskType
}

// 查询是否存在进行中导入任务
func (s *Service) hasTask(groupId, areaId, taskType string) (bool, string) {
	cacheKey := s.cacheKey(groupId, areaId, taskType)
	taskId, err := s.redis.Get(s.ctx, cacheKey).Result()
	if err != nil {
		return false, ""
	}
	progress := s.getProgress(taskId)
	rate := fmt.Sprintf("%v", progress["rate"])
	rateInt, _ := strconv.ParseFloat(rate, 64)
	if rateInt == 0 || rateInt >= 100 {
		return false, ""
	}
	return true, taskId
}

// 清除任务
func (s *Service) clearTask(groupId, areaId, taskType string) bool {
	cacheKey := s.cacheKey(groupId, areaId, taskType)
	taskId, err := s.redis.Get(s.ctx, cacheKey).Result()
	if err != nil {
		return true
	}
	s.redis.Del(s.ctx, cacheKey)
	s.redis.Del(s.ctx, taskId)

	return true
}

// 获取处理进度
func (s *Service) getProgress(taskId string) map[string]interface{} {
	count, err := s.redis.HGet(s.ctx, taskId, "count").Result()
	if err != nil {
		count = "0"
	}
	countInt, _ := strconv.Atoi(count)
	successCount, err := s.redis.HGet(s.ctx, taskId, "successCount").Result()
	if err != nil {
		successCount = "0"
	}
	successCountInt, _ := strconv.Atoi(successCount)
	failCount, err := s.redis.HGet(s.ctx, taskId, "failCount").Result()
	if err != nil {
		failCount = "0"
	}
	failCountInt, _ := strconv.Atoi(failCount)
	doneCount := successCountInt + failCountInt

	var rate int
	if countInt == 0 {
		rate = 0
	} else {
		rateString := fmt.Sprintf("%.2f", float64(doneCount)/float64(countInt))
		rateFloat, _ := strconv.ParseFloat(rateString, 64)
		rate = int(rateFloat * 100)
	}

	filePath, err := s.redis.HGet(s.ctx, taskId, "filePath").Result()
	if err != nil {
		filePath = ""
	}

	isDone, err := s.redis.HGet(s.ctx, taskId, "isDone").Result()
	if err != nil {
		isDone = "0"
	}
	isDoneInt, _ := strconv.Atoi(isDone)

	return map[string]interface{}{
		"count":        countInt,
		"successCount": successCountInt,
		"failCount":    failCountInt,
		"doneCount":    doneCount,
		"isDone":       isDoneInt,
		"rate":         rate,
		"filePath":     filePath,
	}
}

// 构建请求数据
func (s *Service) buildApiDataMap(form *ImportForm, row []string) map[string]interface{} {
	dataMap := map[string]interface{}{}
	mapping := strings.Split(form.Mapping, ",")
	//for k, v := range row {
	//	dataMap[strings.TrimSpace(mapping[k])] = strings.TrimSpace(v)
	//}

	for k, v := range mapping {
		rowLen := len(row)
		if rowLen >= k+1 {
			dataMap[v] = strings.TrimSpace(row[k])
		} else {
			dataMap[v] = ""
		}
	}

	for k, v := range form.ApiParams {
		dataMap[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return dataMap
}

// 调用api接口创建数据
func (s *Service) requestApi(form *ImportForm, dataMap map[string]interface{}) ([]gjson.Result, error) {
	bodyBytes, err := request.Post(s.getApiUrl(form.ApiHost, form.ApiPath), dataMap)
	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(bodyBytes)
	code := gjson.Get(r.Raw, "code").String()
	if code != "200" {
		errMsg := gjson.Get(r.Raw, "msg").String()
		return nil, errors.New(errMsg)
	}
	res := gjson.Get(r.Raw, "data").Array()

	return res, nil
}

// 获取远端api的host地址
func (s *Service) getApiHost(host string) interface{} {
	b, _ := json.Marshal(s.cfg.ApiHost)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m[host]
}

// 获取远端api的完整地址
func (s *Service) getApiUrl(apiHost, apiPath string) string {
	host := s.getApiHost(apiHost)
	return fmt.Sprintf("%v%v", host, apiPath)
}

func (s *Service) log(f func() (string, [][]string, string, error), form *ImportForm) (string, [][]string, string, interface{}, error) {
	collection := s.mongodb.Database("cst_ucenter").Collection("import_excel_log")

	res, err := collection.InsertOne(s.ctx, bson.D{
		{"form", form},
		{"start_time", time.Now().Unix()},
		{"date", time.Now().In(time.Local).Format("2006-01-02 15:04:05")},
		{"status", 0},
	})
	if err != nil {
		return "", nil, "", "", err
	}
	id := res.InsertedID
	filePath, rows, taskId, err := f()
	if err != nil {
		collection.UpdateByID(s.ctx, id,
			bson.D{{"$set", bson.D{
				{"task_id", taskId},
				{"file_path", filePath},
				{"err", err.Error()},
			}}})
	} else {
		collection.UpdateByID(s.ctx, id,
			bson.D{{"$set", bson.D{
				{"task_id", taskId},
				{"file_path", filePath},
			}}})
	}
	return filePath, rows, taskId, id, err
}

func (s *Service) updateLog(id interface{}, update bson.D) (*mongo.UpdateResult, error) {
	collection := s.mongodb.Database("cst_ucenter").Collection("import_excel_log")
	return collection.UpdateByID(s.ctx, id, update)
}

// 导出excel
func (s *Service) exportExcel(form *ExportForm) (string, error) {
	titles := strings.Split(form.HeaderTitle, ",")
	titleKeys := strings.Split(form.HeaderKey, ",")
	if len(titles) != len(titleKeys) {
		return "", errors.New("传参Excel表头参数格式错误")
	}

	taskId := utils.RandStringBytes(10)
	go s.generateExcelFile(form, taskId)

	return taskId, nil
}

// 生成导出excel文件
func (s *Service) generateExcelFile(form *ExportForm, taskId string) error {

	newFile := excelize.NewFile()
	index := newFile.NewSheet("Sheet1")

	titles := strings.Split(form.HeaderTitle, ",")
	for col, value := range titles {
		axis := s.getAxis(col+1, 1)
		newFile.SetCellValue("Sheet1", axis, value)
	}

	_, totalCount, err := s.getTotalPageAndTotalCount(form)
	if err != nil {
		return err
	}

	cacheKey := s.cacheKey(form.ApiParams["group_id"], form.ApiParams["area_id"], form.TaskType)
	s.redis.Set(s.ctx, cacheKey, taskId, time.Second*3600*2)
	s.redis.Del(s.ctx, taskId)
	s.redis.HSet(s.ctx, taskId, "count", totalCount)
	s.redis.Expire(s.ctx, taskId, time.Second*3600*2)

	resChan := make(chan gjson.Result, totalCount)
	go s.exportWorker(form, resChan)

	var page = 0
	for {
		select {
		case ch, ok := <-resChan:
			if !ok {
				filePath, _ := s.getFilePath(form.TaskType)
				if err := newFile.SaveAs(filePath); err != nil {
					return err
				}
				basePath := strings.TrimLeft(filePath, "/www/runtime/go-excel/")
				s.redis.HSet(s.ctx, taskId, "filePath", fmt.Sprintf("%v:%v/%v", s.cfg.Intranet.Ip, s.cfg.Application.Port, basePath))
				s.redis.HSet(s.ctx, taskId, "isDone", 1)
				return nil
			}

			exportData, err := s.buildExportData(form, ch)
			if err != nil {
				return err
			}

			startRow := page*int(pageSize) + 2
			for row, data := range exportData {
				for col, value := range data {
					axis := s.getAxis(col+1, row+startRow)
					newFile.SetCellValue("Sheet1", axis, value)
				}
				s.redis.HIncrBy(s.ctx, taskId, "successCount", 1)
			}
			newFile.SetActiveSheet(index)
			page++
		}
	}

	return nil
}

// 导出worker
func (s *Service) exportWorker(form *ExportForm, resChan chan<- gjson.Result) error {
	totalPage, _, err := s.getTotalPageAndTotalCount(form)
	if err != nil {
		return err
	}

	for i := 1; i <= int(totalPage); i++ {
		res, err := s.getExportDataByPagination(form, int64(i), pageSize)
		if err != nil {
			return err
		}
		resChan <- res
		time.Sleep(time.Millisecond * 50)
		if i == int(totalPage) {
			close(resChan)
		}
	}

	return nil
}

func (s *Service) getExportUrl(url string) string {
	host := strings.Split(url, "/")
	apiHost := s.getApiHost(host[0])
	apiHostString := fmt.Sprintf("%v", apiHost)
	return strings.Replace(url, host[0], apiHostString, 1)
}

// 获取数据总页数
func (s *Service) getTotalPageAndTotalCount(form *ExportForm) (int64, int64, error) {
	res, err := s.getExportDataByPagination(form, 1, pageSize)
	if err != nil {
		return 0, 0, err
	}
	page := gjson.Get(res.Raw, "data.page")
	totalPage := page.Get("totalPage").Int()
	totalCount := page.Get("totalCount").Int()
	return totalPage, totalCount, nil
}

// 分页获取远端要导出的数据
func (s *Service) getExportDataByPagination(form *ExportForm, page, pageSize int64) (gjson.Result, error) {
	apiPath := s.getExportUrl(form.ApiPath)
	query := url.Values{}
	for k, v := range form.ApiParams {
		query.Add(k, fmt.Sprintf("%v", v))
	}
	query.Add("page", fmt.Sprintf("%v", page))
	query.Add("pageSize", fmt.Sprintf("%v", pageSize))

	url := apiPath + "?" + query.Encode()
	bodyBytes, err := request.Get(url)
	if err != nil {
		return gjson.Result{}, err
	}

	r := gjson.ParseBytes(bodyBytes)
	code := gjson.Get(r.Raw, "code").String()
	if code != "200" {
		errMsg := gjson.Get(r.Raw, "msg").String()
		return gjson.Result{}, errors.New(errMsg)
	}
	return r, nil
}

// 构建导出数据
func (s *Service) buildExportData(form *ExportForm, res gjson.Result) ([][]string, error) {
	data := gjson.Get(res.Raw, "data.list").Array()

	var rows [][]string
	for _, v := range data {

		m, ok := gjson.Parse(v.Raw).Value().(map[string]interface{})
		if !ok {
			return nil, errors.New("导出数据格式错误")
		}

		var row []string
		keys := strings.Split(form.HeaderKey, ",")
		for _, v2 := range keys {
			value := fmt.Sprintf("%v", m[v2])
			row = append(row, value)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// 获取导出文件路径
func (s *Service) getFilePath(taskType string) (string, error) {
	exportDir := fmt.Sprintf("%v/%v/%v", "/www/runtime/go-excel/static/export", time.Now().Format("20060102"), taskType)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%v/%v%v.%v", exportDir, time.Now().Unix(), utils.RandStringBytes(5), "xlsx")
	return filePath, nil
}

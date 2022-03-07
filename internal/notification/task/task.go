package task

import (
	"context"
	"cst/internal/notification/admin"
	"cst/internal/notification/area_admin"
	"cst/internal/notification/entity"
	"cst/internal/notification/notify"
	"cst/internal/notification/target"
	"cst/internal/pkg/config"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type task struct {
	db      *gorm.DB
	mongodb *mongo.Client
	logger  *logrus.Logger
	cfg     *config.Config
	mu      sync.RWMutex
	ctx     context.Context
}

func NewTask(db *gorm.DB, mongodb *mongo.Client, logger *logrus.Logger, cfg *config.Config) *task {
	return &task{
		db:      db,
		mongodb: mongodb,
		logger:  logger,
		cfg:     cfg,
	}
}

//SendNotification 发送系统消息通知
func (t *task) SendNotification(doneChan chan int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	notifyRepository := notify.NewRepository(t.db, t.logger)
	notifies, err := notifyRepository.FindAll("type = 0 and receive_type in (0, 1) and send_type = 0 and send_status = 0 and send_time <= ?", time.Now().Unix())
	if err != nil {
		return err
	}
	if len(notifies) > 0 {
		f := func(doneChan chan int) ([]entity.Notify, error) {

			targetRepository := target.NewRepository(t.db, t.logger)
			for _, notify := range notifies {
				admins, err := t.getGroupAdmins(notify)
				if err != nil {
					t.logger.WithFields(logrus.Fields{
						"notify_id": notify.ID,
					}).Error(err.Error())
					return nil, err
				}
				for _, admin := range admins {
					targetRepository.Create(&entity.NotifyTarget{
						NotifyID:  notify.ID,
						Scene:     0,
						AdminID:   admin.ID,
						CreatedAt: time.Now().Unix(),
					})
				}
				areaAdmins, err := t.getAreaAdmins(notify)
				if err != nil {
					t.logger.WithFields(logrus.Fields{
						"notify_id": notify.ID,
					}).Error(err.Error())
					return nil, err
				}
				for _, admin := range areaAdmins {
					targetRepository.Create(&entity.NotifyTarget{
						NotifyID:  notify.ID,
						Scene:     1,
						AdminID:   admin.ID,
						CreatedAt: time.Now().Unix(),
					})
				}
				t.Done(notify.ID)
			}
			return notifies, nil
		}
		t.log(doneChan, f)
	}

	return nil
}

func (t *task) log(doneChan chan int, f func(doneChan chan int) ([]entity.Notify, error)) error {
	collection := t.mongodb.Database("cst_ucenter").Collection("notification_task")

	log, err := collection.InsertOne(t.ctx, bson.D{
		{"start_time", time.Now().Unix()},
		{"date", time.Now().In(time.Local).Format("2006-01-02 15:04:05")},
		{"status", 0},
	})
	id := log.InsertedID
	if err != nil {
		return nil
	}

	res, err := f(doneChan)
	if err != nil {
		collection.UpdateByID(t.ctx, id,
			bson.D{{"$set", bson.D{
				{"status", 2},
				{"end_time", time.Now().Unix()},
				{"err", err.Error()},
				{"res_msg", "执行失败"},
			}}})
	} else {
		collection.UpdateByID(t.ctx, id,
			bson.D{{"$set", bson.D{
				{"status", 1},
				{"end_time", time.Now().Unix()},
				{"data", res},
				{"res_msg", "执行成功"},
			}}})
	}
	return nil
}

//Done 发送完成，更新任务状态
func (t *task) Done(id int) (*entity.Notify, error) {
	notifyRepository := notify.NewRepository(t.db, t.logger)
	return notifyRepository.Update(id, &entity.Notify{SendStatus: 1, CompletionTime: time.Now().Unix()})
}

// 获取通知的集团管理员
func (t *task) getGroupAdmins(notify entity.Notify) ([]entity.Admin, error) {
	var admins []entity.Admin
	if notify.ReceiveType == 0 {
		admin := admin.NewRepository(t.db, t.logger)
		admins, _ = admin.FindAll("is_del = 0 and account_status = 0")
	} else {
		versions, err := notify.GetVersions(t.db)
		if err != nil {
			return nil, err
		}
		for _, version := range versions {
			groupAdmins, _ := version.GetGroupAdmins(t.db, version.VersionID)
			for _, admin := range groupAdmins {
				admins = append(admins, admin)
			}
		}
	}

	return admins, nil
}

// 获取通知的门店管理员
func (t *task) getAreaAdmins(notify entity.Notify) ([]entity.AreaAdmin, error) {
	var admins []entity.AreaAdmin
	if notify.ReceiveType == 0 {
		admin := area_admin.NewRepository(t.db, t.logger)
		admins, _ = admin.FindAll("is_del = 0 and account_status = 0")
	} else {
		versions, err := notify.GetVersions(t.db)
		if err != nil {
			return nil, err
		}
		for _, version := range versions {
			areaAdmins, _ := version.GetAreaAdmins(t.db, version.VersionID)
			for _, admin := range areaAdmins {
				admins = append(admins, admin)
			}
		}
	}

	return admins, nil
}

// 获取所有的集团版本
func (t *task) getGroupVersions() ([]gjson.Result, error) {
	resp, err := http.Get(fmt.Sprintf("%s/group/common/version-index?pagination=0", t.cfg.Versions.Url))
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := gjson.ParseBytes(bodyBytes)
	code := gjson.Get(r.Raw, "code").String()
	errMsg := gjson.Get(r.Raw, "msg").String()
	if code != "200" {
		return nil, fmt.Errorf("请求失败：%s", errMsg)
	}
	return gjson.Get(r.Raw, "data.items").Array(), nil
}

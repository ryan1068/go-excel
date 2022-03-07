package excel

type ImportForm struct {
	FilePath  string            `form:"filePath" binding:"required"`
	TaskType  string            `form:"taskType" binding:"required"`
	Mapping   string            `form:"mapping" binding:"required"`
	ApiHost   string            `form:"apiHost" binding:"required"`
	ApiPath   string            `form:"apiPath" binding:"required"`
	ApiParams map[string]string `form:"apiParams"`
}

type TaskForm struct {
	TaskType string `form:"task_type" binding:"required"`
	GroupId  string `form:"group_id"`
	AreaId   string `form:"area_id"`
}

type ExportForm struct {
	TaskType    string            `form:"taskType" binding:"required"`
	HeaderTitle string            `form:"headerTitle" binding:"required"`
	HeaderKey   string            `form:"headerKey" binding:"required"`
	ApiPath     string            `form:"apiPath" binding:"required"`
	ApiParams   map[string]string `form:"apiParams"`
}

package excel

import (
	"context"
	"cst/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Routers struct {
	service *Service
}

func New(cfg *config.Config, redis *redis.Client, mongodb *mongo.Client, ctx context.Context) Routers {
	return Routers{
		service: &Service{cfg: cfg, redis: redis, mongodb: mongodb, ctx: ctx},
	}
}

func (r Routers) Register(router *gin.RouterGroup) {
	router.GET("/progress", r.View)
	router.GET("/has-task", r.HasTask)
	router.GET("/clear-task", r.ClearTask)
	router.POST("/import", r.Import)
	router.POST("/export", r.Export)
}

func (r Routers) Import(c *gin.Context) {
	form := &ImportForm{}
	if err := c.ShouldBind(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  FirstError(err),
			"data": gin.H{},
		})
		return
	}

	form.ApiParams, _ = c.GetPostFormMap("apiParams")
	taskId, err := r.service.createTask(form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": gin.H{
			"task_id": taskId,
		},
	})
}

func (r Routers) ClearTask(c *gin.Context) {
	taskForm := &TaskForm{}
	if err := c.ShouldBind(taskForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  FirstError(err),
			"data": gin.H{},
		})
		return
	}
	r.service.clearTask(taskForm.GroupId, taskForm.AreaId, taskForm.TaskType)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": gin.H{},
	})
}

func (r Routers) HasTask(c *gin.Context) {
	taskForm := &TaskForm{}
	if err := c.ShouldBind(taskForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  FirstError(err),
			"data": gin.H{},
		})
		return
	}
	hasTask, taskId := r.service.hasTask(taskForm.GroupId, taskForm.AreaId, taskForm.TaskType)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": gin.H{
			"has_task": hasTask,
			"task_id":  taskId,
		},
	})
}

func (r Routers) View(c *gin.Context) {
	taskId := c.Query("task_id")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "任务id不能为空",
			"data": gin.H{},
		})
		return
	}
	progress := r.service.getProgress(taskId)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": progress,
	})
}

func (r Routers) Export(c *gin.Context) {
	form := &ExportForm{}
	if err := c.ShouldBind(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  FirstError(err),
			"data": gin.H{},
		})
		return
	}

	form.ApiParams, _ = c.GetPostFormMap("apiParams")
	taskId, err := r.service.exportExcel(form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": gin.H{
			"task_id": taskId,
		},
	})
}

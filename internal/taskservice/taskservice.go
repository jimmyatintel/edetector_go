package taskservice

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"edetector_go/config"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	work_from_api "edetector_go/internal/work_from_api"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskRequest struct {
	TaskID string `json:"taskID"`
}

type TaskResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

func Start(ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create(config.Viper.GetString("GIN_LOG_FILE"))
	gin.DefaultWriter = io.MultiWriter(f)
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Accept", "Content-Length", "Authorization", "Origin", "X-Requested-With"}
	router.RedirectFixedPath = true
	router.Use(cors.New(corsConfig))
	router.Use(logger.GinLog())
	router.POST("/sendTask", func(c *gin.Context) {
		ReceiveTask(c, ctx)
	})
	router.Run(":5055")
}

func ReceiveTask(c *gin.Context, ctx context.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request format", zap.Any("error", err.Error()))
		return
	}
	handleTaskrequest(ctx, req.TaskID)
	res := TaskResponse{
		IsSuccess: true,
		Message:   "Success",
	}
	c.JSON(http.StatusOK, res)
}

func handleTaskrequest(ctx context.Context, taskid string) {
	logger.Info("Handling task: " + taskid)
	// task_ctx := context.WithValue(ctx, TaskIDKey, taskid)
	message := redis.RedisGetString(taskid)
	content := []byte(message)
	NewPacket := new(packet.TaskPacket)
	err := NewPacket.NewPacket(content)
	if err != nil {
		logger.Error("Error reading task packet:", zap.Any("error", err.Error()), zap.Any("len", len(content)))
		return
	}
	if NewPacket.GetUserTaskType() == "Undefine" {
		nullIndex := bytes.IndexByte(content[76:100], 0)
		logger.Error("Undefine User Task Type: ", zap.String("error", string(content[76:76+nullIndex])))
		return
	}
	logger.Info("Task " + taskid + " " + string(NewPacket.GetUserTaskType()) + " is handling...")
	taskFunc, ok := work_from_api.WorkapiMap[NewPacket.GetUserTaskType()]
	if !ok {
		logger.Error("Function notfound:", zap.Any("name", NewPacket.GetUserTaskType()))
		return
	}
	_, err = taskFunc(NewPacket)
	if err != nil {
		logger.Error(string(NewPacket.GetUserTaskType())+" task failed:", zap.Any("error", err.Error()))
		UsertaskType, ok := task.UserTaskTypeMap[NewPacket.GetUserTaskType()]
		if ok {
			query.Failed_task(NewPacket.GetRkey(), UsertaskType)
		}
		return
	}
	if NewPacket.GetUserTaskType() == "ChangeDetectMode" {
		query.Finish_task(NewPacket.GetRkey(), "ChangeDetectMode")
	}
}

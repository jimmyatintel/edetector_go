package taskservice

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"edetector_go/config"
	"edetector_go/internal/packet"
	work "edetector_go/internal/work"
	work_from_api "edetector_go/internal/work_from_api"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	router.POST("/listscore/:type", func(c *gin.Context) {
		ReceiveUpdateLists(c, ctx)
	})
	router.Run(":5055")
}

func ReceiveTask(c *gin.Context, ctx context.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request format: " + err.Error())
		res := TaskResponse{
			IsSuccess: false,
			Message:   "Invalid request format",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	go handleTaskrequest(ctx, req.TaskID)
	res := TaskResponse{
		IsSuccess: true,
		Message:   "Success",
	}
	c.JSON(http.StatusOK, res)
}

// To-Do (TBD)
func handleTaskrequest(ctx context.Context, taskid string) {
	logger.Info("Handling task: " + taskid)
	// task_ctx := context.WithValue(ctx, TaskIDKey, taskid)
	message := redis.RedisGetString(taskid)
	content := []byte(message)
	NewPacket := new(packet.TaskPacket)
	err := NewPacket.NewPacket(content)
	if err != nil {
		logger.Error("Error reading task packet: " + err.Error())
		query.Update_task_status_by_taskid(taskid, 6)
		return
	}
	if NewPacket.GetUserTaskType() == "Undefine" {
		nullIndex := bytes.IndexByte(content[76:100], 0)
		logger.Error("Undefine User Task Type: " + string(content[76:76+nullIndex]))
		query.Update_task_status_by_taskid(taskid, 6)
		return
	}
	logger.Info("Task " + taskid + " " + string(NewPacket.GetUserTaskType()) + " is handling...")
	taskFunc, ok := work_from_api.WorkapiMap[NewPacket.GetUserTaskType()]
	if !ok {
		logger.Error("Function notfound:" + string(NewPacket.GetUserTaskType()))
		query.Update_task_status_by_taskid(taskid, 6)
		return
	}
	_, err = taskFunc(NewPacket)
	if err != nil {
		logger.Error("Task " + string(NewPacket.GetUserTaskType()) + " failed: " + err.Error())
		query.Update_task_status_by_taskid(taskid, 6)
		return
	}
	if NewPacket.GetUserTaskType() == "ChangeDetectMode" {
		query.Finish_task(NewPacket.GetRkey(), "ChangeDetectMode")
	}
}

func ReceiveUpdateLists(c *gin.Context, ctx context.Context) {
	listType := c.Param("type")
	redis.RedisSet("update_"+listType+"_list", "0")
	var req interface{}

	if listType == "hack" {
		var hackReq work.HackReq
		if err := c.ShouldBindJSON(&hackReq); err != nil {
			ErrorResponse(c, err, "Invalid request format")
			return
		}
		req = hackReq
	} else if listType == "white" || listType == "black" {
		var wbReq work.WhiteBlackReq
		if err := c.ShouldBindJSON(&wbReq); err != nil {
			ErrorResponse(c, err, "Invalid request format")
			return
		}
		req = wbReq
	} else {
		ErrorResponse(c, nil, "Invalid list type")
		return
	}
	res := TaskResponse{
		IsSuccess: true,
		Message:   "Success",
	}
	c.JSON(http.StatusOK, res)
	go work.UpdateLists(listType, req)
}

func ErrorResponse(c *gin.Context, err error, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": false,
		"message":   msg + ": " + err.Error(),
	})
	logger.Error(msg + ": " + err.Error())
}

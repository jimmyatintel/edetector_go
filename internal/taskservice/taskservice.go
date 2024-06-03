package taskservice

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	work "edetector_go/internal/work"
	work_from_api "edetector_go/internal/work_from_api"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	rq "edetector_go/pkg/redis/query"

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
	logger.Info("Task " + req.TaskID + " received")
}

// To-Do (TBD)
func handleTaskrequest(ctx context.Context, taskid string) {
	logger.Info("Handling task: " + taskid)
	query.Update_task_status_by_taskid(taskid, 2)
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
	ttype := NewPacket.GetUserTaskType()
	key := NewPacket.GetRkey()
	if ttype == "Undefine" {
		nullIndex := bytes.IndexByte(content[76:100], 0)
		logger.Error("Undefine User Task Type: " + string(content[76:76+nullIndex]))
		query.Update_task_status_by_taskid(taskid, 6)
		return
	}
	logger.Info("Task " + taskid + " " + string(ttype) + " is handling...")
	if ttype == "StartRemove" && rq.GetStatus(key) == 0 {
		DeleteAgentData(key)
		return
	}
	taskFunc, ok := work_from_api.WorkapiMap[ttype]
	if !ok {
		logger.Error("Function notfound:" + string(ttype))
		query.Update_task_status_by_taskid(taskid, 6)
		return
	}
	_, err = taskFunc(NewPacket)
	if err != nil {
		logger.Error("Task " + string(ttype) + " failed: " + err.Error())
		query.Update_task_status_by_taskid(taskid, 6)
		return
	}
	if ttype == "ChangeDetectMode" {
		query.Finish_task(key, "ChangeDetectMode")
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

func DeleteAgentData(key string) {
	query.DeleteAgent(key)
	redisData := redis.GetKeysMatchingPattern(key + "*")
	for _, r := range redisData {
		err := redis.RedisDelete(r)
		if err != nil {
			logger.Error("Error deleting data from redis: " + err.Error())
		}
	}
}

func RetryTask(key string, tasktype string, retryTask task.TaskType) error {
	// get retry count from radis
	retryCount, err := redis.RedisGetInt(key + "-RetryCount")
	if err != nil {
		logger.Error("Get retry count failed: " + err.Error())
		query.Failed_task(key, tasktype, 7)
		// reset retry count
		redisErr := redis.RedisSet(key+"-RetryCount", 0)
		if redisErr != nil {
			logger.Error("Reset retry count failed: " + redisErr.Error())
			return redisErr
		}
		return err
	}
	if retryCount >= config.Viper.GetInt("RETRY_COUNT") {
		query.Failed_task(key, tasktype, 7)
		// reset retry count
		redisErr := redis.RedisSet(key+"-RetryCount", 0)
		if redisErr != nil {
			logger.Error("Reset retry count failed: " + redisErr.Error())
			return redisErr
		}
		return errors.New("retry count is over")
	}

	// get agent ip and mac from mariaDB
	ip := query.GetMachineIP(key)
	mac := query.GetMachineMAC(key)

	// send ResendCollect task to agent
	err = clientsearchsend.SendUserTCPtoClientWithoutP(ip, mac, retryTask, "")
	if err != nil {
		logger.Error("Send ResendCollect task failed: " + err.Error())
		return err
	}

	// update retry count
	err = redis.RedisSet_AddInteger((key + "-RetryCount"), 1)
	if err != nil {
		logger.Error("Update retry count failed: " + err.Error())
		return err
	}

	return nil
}

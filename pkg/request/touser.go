package request

import (
	"bytes"
	"context"
	"edetector_go/config"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Request struct {
	DeviceId string `json:"deviceId"`
}

type ReadyRequest struct {
	TaskId   string `json:"taskId"`
	DllPaths string `json:"dllPaths"`
}

type ReadyData struct {
	TaskId   string
	DllPaths string
	Failed   string
	RedisKey string
	Progress int
	LoadDll  bool
}

func RequestToUser(id string) {
	request := Request{
		DeviceId: id,
	}
	// Marshal payload into JSON
	payload, err := json.Marshal(request)
	if err != nil {
		logger.Error("Error marshaling JSON: " + err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create an HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	ip := config.Viper.GetString("WS_HOST")
	port := config.Viper.GetString("WS_PORT")
	path := fmt.Sprintf("http://%s:%s/updateTask", ip, port)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("Error creating HTTP request: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the HTTP request
	response, err := client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				logger.Error("Request timed out: " + err.Error())
				return
			}
		default:
			logger.Error("Error sending HTTP request: " + err.Error())
			return
		}
	}
	defer response.Body.Close()
	// Check the response status code
	if response.StatusCode != http.StatusOK {
		logger.Error("Request failed with status code: " + fmt.Sprint(response.StatusCode))
		return
	}
}

// LoadDumpReady updates progress in redis and informs API
func LoadDumpReady(info ReadyData) {
	// if progress is the same --> no need to update
	if !redis.NeedToUpdateProgress(info.TaskId, info.Progress) {
		return
	}

	// check taskId exists in pendingDump:USERID
	if !info.LoadDll && !redis.CheckDumpTaskExists(info.TaskId) {
		logger.Warn("TaskId does not exist in pendingDump:USERID")
		return
	}
	logger.Debug("TaskId exists in pendingDump:USERID")

	// Marshal payload into JSON
	payload, err := json.Marshal(ReadyRequest{
		TaskId:   info.TaskId,
		DllPaths: info.DllPaths,
	})
	if err != nil {
		logger.Error("Error marshaling JSON: " + err.Error())
		return
	}

	// update info in redis
	redis.UpdateDumpTaskInfo(info.TaskId, info.Failed, info.Progress)
	logger.Debug("Updated task info in redis")

	// update pending dump in redis if the task is not loadDll
	if !info.LoadDll {
		if info.Progress == 100 {
			redis.UpdatePendingDump(info.TaskId, 1)
		} else if info.Progress == -1 {
			redis.UpdatePendingDump(info.TaskId, -1)
		} else if info.Progress == -2 {
			redis.UpdatePendingDump(info.TaskId, -2)
		}
	}
	logger.Debug("Updated pending dump in redis")

	// remove key from redis if the progress is 100, -1, -2
	if info.Progress == 100 || info.Progress == -1 || info.Progress == -2 {
		if err := redis.RedisDelete(info.RedisKey); err != nil {
			logger.Error("Error deleting key from redis: " + err.Error())
		}
	}
	logger.Debug("Deleted task key from redis (progress = 100, -1, -2)")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create an HTTP request
	ip := config.Viper.GetString("WS_HOST")
	port := config.Viper.GetString("WS_PORT")
	path := fmt.Sprintf("http://%s:%s/loadDumpReady", ip, port)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("Error creating HTTP request: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				logger.Error("Request timed out: " + err.Error())
				return
			}
		default:
			logger.Error("Error sending HTTP request: " + err.Error())
			return
		}
	}
	defer response.Body.Close()
	logger.Debug("Sent request to API")

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		logger.Error("Request failed with status code: " + fmt.Sprint(response.StatusCode))
		return
	}
	logger.Debug("Receive response from API")
}

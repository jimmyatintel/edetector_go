package redis

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/request"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type ClientOnlineStatus struct {
	Status int
	Time   string
}

func (c *ClientOnlineStatus) Marshal() string {
	json, err := json.Marshal(c)
	if err != nil {
		logger.Error("Error in json marshal" + err.Error())
	}
	return string(json)
}

func Online(KeyNum string) {
	currentTime := time.Now().Format(time.RFC3339)
	onlineStatusInfo := ClientOnlineStatus{
		Status: 1,
		Time:   currentTime,
	}
	err := RedisSet(KeyNum, onlineStatusInfo.Marshal())
	if err != nil {
		logger.Error("Update online failed:" + err.Error())
		return
	}
}

func Offline(KeyNum string) {
	currentTime := time.Now().Format(time.RFC3339)
	onlineStatusInfo := ClientOnlineStatus{
		Status: 0,
		Time:   currentTime,
	}
	err := RedisSet(KeyNum, onlineStatusInfo.Marshal())
	if err != nil {
		logger.Error("Update offline failed:" + err.Error())
		return
	}
	handlingTasks, err := query.Load_stored_task("nil", KeyNum, 2, "nil")
	if err != nil {
		logger.Error("Get handling tasks failed:", zap.Any("error", err.Error()))
		return
	}
	for _, t := range handlingTasks {
		query.Failed_task(KeyNum, t[3])
	}
	logger.Info("Offline and let all task failed: ", zap.Any("message", KeyNum))
	request.RequestToUser(KeyNum)
}

func GetStatus(content string) int {
	var clientStatus ClientOnlineStatus
	err := json.Unmarshal([]byte(content), &clientStatus)
	if err != nil {
		logger.Error("Get status failed:", zap.Any("error", err.Error()))
		return -1
	}
	return clientStatus.Status
}

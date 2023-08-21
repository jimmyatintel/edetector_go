package redis

import (
	"edetector_go/pkg/logger"
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

func Online(KeyNum string) error {
	currentTime := time.Now().Format(time.RFC3339)
	onlineStatusInfo := ClientOnlineStatus{
		Status: 1,
		Time:   currentTime,
	}
	return RedisSet(KeyNum, onlineStatusInfo.Marshal())
}

func Offline(KeyNum string) error {
	currentTime := time.Now().Format(time.RFC3339)
	onlineStatusInfo := ClientOnlineStatus{
		Status: 0,
		Time:   currentTime,
	}
	return RedisSet(KeyNum, onlineStatusInfo.Marshal())
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

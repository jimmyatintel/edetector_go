package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"encoding/json"
	"time"
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
	err := redis.RedisSet(KeyNum, onlineStatusInfo.Marshal())
	if err != nil {
		logger.Error("Update online failed:" + err.Error())
		return
	}
}

func Offline(KeyNum string) {
	if GetStatus(KeyNum) == 0 { // already offline
		return
	}
	redis.RedisSet_AddInteger("OnlineClientCount", -1)
	currentTime := time.Now().Format(time.RFC3339)
	onlineStatusInfo := ClientOnlineStatus{
		Status: 0,
		Time:   currentTime,
	}
	err := redis.RedisSet(KeyNum, onlineStatusInfo.Marshal())
	if err != nil {
		logger.Error("Update offline failed:" + err.Error())
		return
	}
	handlingTasks1, err := query.Load_stored_task("nil", KeyNum, 1, "nil")
	if err != nil {
		logger.Error("Get handling tasks1 failed: " + err.Error())
		return
	}
	handlingTasks2, err := query.Load_stored_task("nil", KeyNum, 2, "nil")
	if err != nil {
		logger.Error("Get handling tasks2 failed: " + err.Error())
		return
	}
	for _, t := range handlingTasks1 {
		query.Failed_task(KeyNum, t[3], 7)
	}
	for _, t := range handlingTasks2 {
		query.Failed_task(KeyNum, t[3], 7)
	}
	// request.RequestToUser(KeyNum)
}

func GetStatus(KeyNum string) int {
	content := redis.RedisGetString(KeyNum)
	var clientStatus ClientOnlineStatus
	err := json.Unmarshal([]byte(content), &clientStatus)
	if err != nil {
		logger.Error("Get status failed: " + err.Error())
		return -1
	}
	return clientStatus.Status
}

func GetTime(KeyNum string) string {
	content := redis.RedisGetString(KeyNum)
	var clientStatus ClientOnlineStatus
	err := json.Unmarshal([]byte(content), &clientStatus)
	if err != nil {
		logger.Error("Get online time failed: " + err.Error())
		return "1970-01-01T00:00:00Z"
	}
	return clientStatus.Time
}

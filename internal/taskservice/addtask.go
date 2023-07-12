package taskservice

import (
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/redis"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"encoding/json"
	"github.com/google/uuid"
)

func AddTask2db(deviceId string, work task.UserTaskType, msg string) error {
	taskId := uuid.NewString()

	// store into mariaDB
	query := "INSERT INTO task (task_id, client_id, type, status, timestamp) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)"
	_, err := mariadb.DB.Exec(query, taskId, deviceId, work, 0)
	if err != nil {
		return err
	}
	// store into redis
	var pkt = packet.TaskPacket{
		Key: deviceId,
		Work: work,
		User: "1",
		Message: msg,
	}
	pktString, err := json.Marshal(pkt)
	if err != nil {
		return err
	}

	redis.Redis_set(taskId, string(pktString))
	return nil
}
package taskservice

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
)

type task_columns struct {
	clientid string
	taskid   string
	status   int
}

func loadfromdb(q task_columns) []task_columns {
	result := query.Load_stored_task(q.taskid, q.clientid, q.status)
	var ret []task_columns
	for _, v := range result {
		tmp := task_columns{}
		tmp.clientid = v[1]
		tmp.taskid = v[0]
		tmp.status, _ = strconv.Atoi(v[2])
		ret = append(ret, tmp)
	}
	return ret
}

func loadallunhandletask() []task_columns {
	var q task_columns
	q.taskid = "nil"
	q.clientid = "nil"
	q.status = 0
	return loadfromdb(q)
}

func loadallprocesstask() []task_columns {
	var q task_columns
	q.taskid = "nil"
	q.clientid = "nil"
	q.status = 1
	return loadfromdb(q)
}

func AddTask(deviceId string, work task.UserTaskType, msg string) error {
	taskId := uuid.NewString()
	// store into mariaDB
	query := "INSERT INTO task (task_id, client_id, type, status, progress, timestamp) VALUES (?, ?, ?, ?, 0, CURRENT_TIMESTAMP)"
	_, err := mariadb.DB.Exec(query, taskId, deviceId, work, 0)
	if err != nil {
		return err
	}
	// store into redis
	var pkt = packet.TaskPacket{
		Key:     deviceId,
		Work:    work,
		User:    "1",
		Message: msg,
	}
	pktString, err := json.Marshal(pkt)
	if err != nil {
		return err
	}
	err = redis.Redis_set(taskId, string(pktString))
	if err != nil {
		return err
	}
	return nil
}

func Finish_task(clientid string, tasktype string) {
	taskID := Find_task_id(clientid, tasktype)
	Change_task_status(taskID, 3)
	Change_task_timestamp(clientid, tasktype)
	RequestToUser(clientid)
}

func Find_task_id(clientid string, tasktype string) string {
	return query.Select_task_id(clientid, tasktype)
}

func Change_task_status(taskid string, status int) {
	query.Update_task_status(taskid, status)
	if status == 3 {
		clientid, tasktype := query.Get_task_info(taskid)
		if clientid != "" && tasktype != ""{
			Change_task_timestamp(clientid, tasktype)
		}
	}
}

func Change_task_timestamp(clientid string, tasktype string) {
	query.Update_task_timestamp(clientid, tasktype)
}
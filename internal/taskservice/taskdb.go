package taskservice

import (
	"edetector_go/pkg/mariadb/query"
)

func Finish_task(clientid string, tasktype string) {
	taskid := query.Get_task_id(clientid, tasktype)
	Change_task_status(taskid, 3)
	RequestToUser(clientid)
	if tasktype == "ChangeDetectMode" {
		return
	}
	Change_task_timestamp(clientid, tasktype)
}

func Failed_task(clientid string, tasktype string) {
	taskid := query.Get_task_id(clientid, tasktype)
	Change_task_status(taskid, 4)
	RequestToUser(clientid)
}

func Change_task_status(taskid string, status int) {
	query.Update_task_status(taskid, status)
}

func Change_task_timestamp(clientid string, tasktype string) {
	query.Update_task_timestamp(clientid, tasktype)
}
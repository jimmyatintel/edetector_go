package taskservice

import (
	"edetector_go/pkg/mariadb/query"
	"strconv"
)

type task_columns struct {
	clientid string
	taskid  string
	status  int
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

func Find_task_id(clientid string, tasktype string) string {
	return query.Select_task_id(clientid, tasktype)
}

func Change_task_status(taskid string, status int) {
	query.Update_task_status(taskid, status)
}

func Change_task_timestamp(clientid string, tasktype string) {
	query.Update_task_timestamp(clientid, tasktype)
}

package taskservice

import (
	"edetector_go/pkg/mariadb/query"
	"strconv"
)

type task_columns struct {
	agentid string
	taskid  string
	status  int
}

func loadfromdb(q task_columns) []task_columns {
	result := query.Load_stored_task(q.taskid, q.agentid, q.status)
	var ret []task_columns
	for _, v := range result {
		tmp := task_columns{}
		tmp.agentid = v[1]
		tmp.taskid = v[0]
		tmp.status, _ = strconv.Atoi(v[2])
		ret = append(ret, tmp)
	}
	return ret
}

func loadallunhandletask() []task_columns {
	var q task_columns
	q.status = 0
	return loadfromdb(q)
}

func loadallprocesstask() []task_columns {
	var q task_columns
	q.status = 1
	return loadfromdb(q)
}

func change_task_status(taskid string, agentid string, status int) {
	query.Update_task_status(taskid, status)
}

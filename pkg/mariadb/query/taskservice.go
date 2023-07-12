package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"strconv"
)

func Load_stored_task(taskid string, client_id string, status int) [][]string {
	qu := "select task_id, client_id, status from task where "
	var result [][]string
	if client_id != "nil" {
		qu = qu + "client_id = " + client_id
	}
	if taskid != "nil" {
		qu = qu + "task_id = " + taskid
	}
	if status != -1 {
		qu = qu + "status = " + strconv.Itoa(status)
	}
	res, err := mariadb.DB.Query(qu)
	if err != nil {
		logger.Error(err.Error())
		return result
	}
	defer res.Close()
	l, _ := res.Columns()
	for res.Next() {
		tmp := make([]string, len(l))
		err := res.Scan(&tmp[0], &tmp[1], &tmp[2])
		if err != nil {
			logger.Error(err.Error())
			return result
		}
		result = append(result, tmp)
	}
	return result
}

func Update_task_status(taskid string, status int) {
	_, err := mariadb.DB.Exec("update task set status = ? where task_id = ?", status, taskid)
	if err != nil {
		logger.Error(err.Error())
	}
}

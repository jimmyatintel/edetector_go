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

func Select_task_id(clientid string, tasktype string) string {
	var taskID string
	row := mariadb.DB.QueryRow("SELECT task_id FROM task WHERE client_id = ? AND type = ? AND status = ?", clientid, tasktype, 2)
	err := row.Scan(&taskID)
	if err != nil {
		logger.Error(err.Error())
	}
	return taskID
}

func Update_task_status(taskid string, status int) {
	if status == 1 {
		_, err := mariadb.DB.Exec("update task set status = ? where task_id = ?", status, taskid)
		if err != nil {
			logger.Error(err.Error())
		}
	} else if status == 2 {
		_, err := mariadb.DB.Exec("UPDATE task SET status = ?, start = 1 WHERE task_id = ?", status, taskid)
		if err != nil {
			logger.Error(err.Error())
		}
	} else if status == 3 {
		_, err := mariadb.DB.Exec("UPDATE task SET status = ?, finish = 1 WHERE task_id = ?", status, taskid)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func Update_task_timestamp(clientid string, tasktype string) {
	col := ""
	if tasktype == "StartScan" {
		col = "scan_finish_time"
	} else if tasktype == "StartCollect" {
		col = "collect_finish_time"
	} else if tasktype == "StartGetDrive" {
		col = "file_finish_time"
	}
	qu := "update client_task_status set " + col +" = CURRENT_TIMESTAMP where client_id = ?"
	_, err := mariadb.DB.Exec(qu, clientid)
	if err != nil {
		logger.Error(err.Error())
	}
}
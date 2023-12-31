package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/request"
	"strconv"
)

func Load_stored_task(task_id string, client_id string, status int, tasktype string) ([][]string, error) {
	qu := "SELECT task_id, client_id, status, type FROM task where "
	var result [][]string
	if client_id != "nil" {
		qu = qu + "client_id = \"" + client_id + "\" AND "
	}
	if task_id != "nil" {
		qu = qu + "task_id = \"" + task_id + "\" AND "
	}
	if status != -1 {
		qu = qu + "status = " + strconv.Itoa(status) + " AND "
	}
	if tasktype != "nil" {
		qu = qu + "type = \"" + tasktype + "\" AND "
	}
	qu = qu + "1=1"
	res, err := mariadb.DB.Query(qu)
	if err != nil {
		return result, err
	}
	defer res.Close()
	l, _ := res.Columns()
	for res.Next() {
		tmp := make([]string, len(l))
		err := res.Scan(&tmp[0], &tmp[1], &tmp[2], &tmp[3])
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	return result, nil
}

func Update_progress(progress int, clientid string, tasktype string) {
	qu := "update task set progress = ? where client_id = ? and type = ? and status = 2"
	_, err := mariadb.DB.Exec(qu, progress, clientid, tasktype)
	if err != nil {
		logger.Error("Update failed: " + err.Error())
	}
}

func Update_task_status(clientid string, tasktype string, old_status int, new_status int) int {
	result, err := mariadb.DB.Exec("update task set status = ? where client_id = ? and type = ? and status = ?", new_status, clientid, tasktype, old_status)
	if err != nil {
		logger.Error("Error Update_task_status: " + err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rowsAffected: " + err.Error())
	}
	return int(rowsAffected)
}

func Update_task_timestamp(clientid string, tasktype string) {
	col := ""
	if tasktype == "StartScan" {
		col = "scan_finish_time"
	} else if tasktype == "StartCollect" {
		col = "collect_finish_time"
	} else if tasktype == "StartGetDrive" {
		col = "file_finish_time"
	} else if tasktype == "StartGetImage" {
		col = "image_finish_time"
	} else {
		return
	}
	qu := "update client_task_status set " + col + " = CURRENT_TIMESTAMP where client_id = ?"
	_, err := mariadb.DB.Exec(qu, clientid)
	if err != nil {
		logger.Error("Error Update_task_timestamp: " + err.Error())
	}
}

func Finish_task(clientid string, tasktype string) {
	rowsAffected := Update_task_status(clientid, tasktype, 2, 3)
	if tasktype == "ChangeDetectMode" {
		return
	}
	if rowsAffected > 0 {
		Update_task_timestamp(clientid, tasktype)
		request.RequestToUser(clientid)
	}
}

func Terminated_task(clientid string, tasktype string) {
	rowsAffected := Update_task_status(clientid, tasktype, 2, 4)
	if rowsAffected > 0 {
		Update_task_timestamp(clientid, tasktype)
		request.RequestToUser(clientid)
	}
}

func Failed_task(clientid string, tasktype string) {
	rowsAffected := Update_task_status(clientid, tasktype, 2, 5)
	if tasktype == "ChangeDetectMode" {
		return
	}
	if rowsAffected > 0 {
		Update_task_timestamp(clientid, tasktype)
		request.RequestToUser(clientid)
	}
}

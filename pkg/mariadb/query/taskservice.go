package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/request"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
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
	fmt.Println(qu)
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
	if len(result) == 0 {
		err := errors.New("task not found")
		return result, err
	}
	return result, nil
}

func Update_task_status(taskid string, status int) {
	_, err := mariadb.DB.Exec("update task set status = ? where task_id = ?", status, taskid)
	if err != nil {
		logger.Error(err.Error())
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
	} else {
		return
	}
	qu := "update client_task_status set " + col + " = CURRENT_TIMESTAMP where client_id = ?"
	_, err := mariadb.DB.Exec(qu, clientid)
	if err != nil {
		logger.Error(err.Error())
	}
}

func Finish_task(clientid string, tasktype string) {
	result, err := Load_stored_task("nil", clientid, -1, tasktype)
	if err != nil {
		logger.Error("Get handling tasks failed:", zap.Any("error", err.Error()))
		return
	}
	taskid := result[0][0]
	Update_task_status(taskid, 3)
	if tasktype == "ChangeDetectMode" {
		return
	}
	Update_task_timestamp(clientid, tasktype)
	request.RequestToUser(clientid)
}

func Terminated_task(taskid string, clientid string, tasktype string) {
	Update_task_status(taskid, 4)
	Update_task_timestamp(clientid, tasktype)
	request.RequestToUser(clientid)
}

func Failed_task(clientid string, tasktype string) {
	result, err := Load_stored_task("nil", clientid, -1, tasktype)
	if err != nil {
		logger.Error("Get handling tasks failed:", zap.Any("error", err.Error()))
		return
	}
	taskid := result[0][0]
	Update_task_status(taskid, 5)
	Update_task_timestamp(clientid, tasktype)
	request.RequestToUser(clientid)
}

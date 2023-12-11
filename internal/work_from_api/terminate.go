package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"os"
	"path/filepath"
)

var disks []string

func init() {
	for i := 'A'; i <= 'Z'; i++ {
		disks = append(disks, string(i))
	}
	disks = append(disks, "Linux")
}

func Terminate(p packet.UserPacket) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("Terminate: " + key + "::" + p.GetMessage())
	handlingTasks, err := query.Load_stored_task("nil", key, 2, "nil")
	if err != nil {
		return task.FAIL, err
	}
	for _, t := range handlingTasks {
		if t[3] == "StartGetDrive" {
			for _, disk := range disks {
				path := filepath.Join("fileUnstage", key+"."+string(disk)+".txt")
				if file.FileExists(path) {
					os.Remove(path)
				}
			}
			query.Terminate_handling_task(key, t[3])
		} else if t[3] == "StartCollect" {
			path := filepath.Join("dbUnstage", key+".db")
			if file.FileExists(path) {
				os.Remove(path)
			}
			query.Terminate_handling_task(key, t[3])
		} else if t[3] != "Terminate" {
			query.Terminated_task(key, t[3], 2)
		}
	}
	err = clientsearchsend.SendUserTCPtoClient(p, task.TERMINATE_ALL, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

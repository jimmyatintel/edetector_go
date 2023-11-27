package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"os"
	"path/filepath"
)

func Terminate(p packet.UserPacket) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("Terminate: " + key + "::" + p.GetMessage())
	handlingTasks, err := query.Load_stored_task("nil", key, 2, "nil")
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-terminateDrive", 0)
	redis.RedisSet(key+"-terminateCollect", 0)
	for _, t := range handlingTasks {
		if t[3] != "Terminate" {
			query.Terminated_task(key, t[3])
		}
		if t[3] == "StartGetDrive" {
			for disk := 'A'; disk <= 'Z'; disk++ {
				path := filepath.Join("fileUnstage", key+"."+string(disk)+".txt")
				if file.FileExists(path) {
					os.Remove(path)
				}
				processingPath := filepath.Join("fileUnstage", key+"."+string(disk)+".txt.processing")
				if file.FileExists(processingPath) {
					redis.RedisSet(key+"-terminateDrive", 1)
				}
			}
			path := filepath.Join("fileUnstage", key+".Linux.txt")
			if file.FileExists(path) {
				os.Remove(path)
			}
			processingPath := filepath.Join("fileUnstage", key+".Linux.txt.processing")
			if file.FileExists(processingPath) {
				redis.RedisSet(key+"-terminateDrive", 1)
			}
		} else if t[3] == "StartCollect" {
			path := filepath.Join("dbUnstage", key+".db")
			if file.FileExists(path) {
				os.Remove(path)
			}
			processingPath := filepath.Join("dbUnstage", key+".db.processing")
			if file.FileExists(processingPath) {
				redis.RedisSet(key+"-terminateCollect", 1)
			}
		}
	}
	err = clientsearchsend.SendUserTCPtoClient(p, task.TERMINATE_ALL, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

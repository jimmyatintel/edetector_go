package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"net"
)

func FinishTerminate(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("FinishTerminate: " + key + "**" + p.GetMessage())
	handlingTasks, err := query.Load_stored_task("nil", key, 2, "nil")
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-terminateFinishIteration", 0)
	redis.RedisSet(key+"-terminateDrive", 0)
	redis.RedisSet(key+"-terminateCollect", 0)
	for _, t := range handlingTasks {
		if t[3] == "StartScan" || t[3] == "StartGetImage" {
			query.Terminated_task(key, t[3])
		} else if t[3] == "StartGetDrive" {
			redis.RedisSet(key+"-terminateDrive", 1)
		} else if t[3] == "StartCollect" {
			redis.RedisSet(key+"-terminateCollect", 1)
		}
	}
	redis.RedisSet(key+"-terminateFinishIteration", 1)
	if redis.RedisGetInt(key+"-terminateDrive") == 0 && redis.RedisGetInt(key+"-terminateCollect") == 0 {
		query.Finish_task(key, "Terminate")
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

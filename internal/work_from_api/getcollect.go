package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
)

func StartCollect(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartCollect: " + p.GetRkey() + "**" + p.GetMessage())
	redis.RedisSet(p.GetRkey()+"-CollectProgress", 0)
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_COLLECT_INFO, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

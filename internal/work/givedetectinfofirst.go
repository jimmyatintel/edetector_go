package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"net"

	"go.uber.org/zap"
)

func GiveDetectInfoFirst(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	// front process back netowork
	logger.Info("GiveDetectInfoFirst: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	redis.RedisSet(key+"-DetectMsg", "")
	rt := query.First_detect_info(p.GetRkey(), p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.UPDATE_DETECT_MODE, rt, conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDetectInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	return task.SUCCESS, nil
}

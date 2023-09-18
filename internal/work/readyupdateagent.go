package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"strings"

	"net"

	"go.uber.org/zap"
)

func ReadyUpdateAgent(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("ReadyUpdateAgent: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	// update progress
	if strings.Split(p.GetMessage(), "/")[0] == "1" {
		collectFirstPart = float64(config.Viper.GetInt("COLLECT_FIRST_PART"))
		collectSecondPart = 100 - collectFirstPart
		go updateCollectProgress(key)
	}
	progress, err := getProgressByMsg(p.GetMessage(), collectFirstPart)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-CollectProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

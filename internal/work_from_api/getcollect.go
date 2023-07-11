package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"

	"go.uber.org/zap"
)

func StartCollection(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartCollection: ", zap.Any("message", p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_COLLECT_INFO_DATA, p.GetMessage(), "undefine")
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

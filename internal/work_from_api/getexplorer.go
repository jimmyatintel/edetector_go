package workfromapi

import (
	"edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"

	"go.uber.org/zap"
)

func StartGetExplorer(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartGetExplorer: ", zap.Any("message", p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.EXPLORER_INFO, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
    return task.SUCCESS, nil
}

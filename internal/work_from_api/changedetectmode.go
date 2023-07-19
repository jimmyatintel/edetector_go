package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"

	"go.uber.org/zap"
)

func ChangeDetectMode(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("ChangeDetectMode: ", zap.Any("message", p.GetRkey() + ", Msg: " + p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.UPDATE_DETECT_MODE, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

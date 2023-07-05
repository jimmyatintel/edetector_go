package workfromapi

import (
	"edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func ChangeDetectMode(p packet.UserPacket, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("ChangeDetectMode: ", zap.Any("message", p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.UPDATE_DETECT_MODE, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
    return task.SUCCESS, nil
}

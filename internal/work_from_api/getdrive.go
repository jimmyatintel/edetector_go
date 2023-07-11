package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	// work "edetector_go/internal/work"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func StartGetDrive(p packet.UserPacket, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("StartGetDrive: ", zap.Any("message", p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_DRIVE, p.GetMessage(), "worker")
	if err != nil {
		return task.FAIL, err
	}
    return task.SUCCESS, nil
}
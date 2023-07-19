package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"

	"go.uber.org/zap"
)

func StartGetDrive(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartGetDrive: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_DRIVE, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

package workfromapi

import (
	"edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func StartGetDrive(p packet.UserPacket, Key *string, conn net.Conn) (task.TaskResult, error) {
	// logger.Info("StartGetDrive: ", zap.Any("message", p.GetMessage()))
	// err := clientsearchsend.SendUserTCPtoClient(p, task.GET_DRIVE, p.GetMessage(), "worker")
	// if err != nil {
	// 	return task.FAIL, err
	// }
    // return task.SUCCESS, nil
	logger.Info("StartGetExplorer: ", zap.Any("message", p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.EXPLORER_INFO, "C|NTFS|Explorer|ScheduleName|0|2048", "worker")
	if err != nil {
		return task.FAIL, err
	}
    return task.SUCCESS, nil
}
package workfromapi

import (
	"edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func StartScan(p packet.UserPacket, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("StartScan: ", zap.Any("message", p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_SCAN_INFO_DATA, p.GetMessage(), "worker")
	if err != nil {
		return task.FAIL, err
	}
	// err_ := clientsearchsend.SendUserTCPtoClient(p, task.GET_SCAN_INFO_DATA, p.GetMessage(), "detect")
	// if err_ != nil {
	// 	return task.FAIL, err
	// }
    return task.SUCCESS, nil
}

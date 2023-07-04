package workfromapi

import (
	_ "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func ChangeDetectMode(p packet.UserPacket, Key *string, conn net.Conn) (task.TaskResult, error) {

	logger.Info("ChangeDetectMode: ", zap.Any("message", p.GetMessage()))

	// Inform agent: "0|0"
	// err := clientsearchsend.SendUserTCPtoClient(p.GetRkey())
	// if err != nil {
	// 	return task.FAIL, err
	// }

	

	return task.SUCCESS, nil
}

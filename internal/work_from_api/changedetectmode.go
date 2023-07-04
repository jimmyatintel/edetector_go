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

	// Inform agent
	logger.Info("GiveDetectProcessOver: ", zap.Any("message", p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.UPDATE_DETECT_MODE,
		Message:    p.GetMessage(), // "0|0"
	}
	err := clientsearchsend.SendUserTCPtoClient(send_packet.Fluent(), p.GetRkey())
	if err != nil {
		return task.FAIL, err
	}

    return task.SUCCESS, nil
}

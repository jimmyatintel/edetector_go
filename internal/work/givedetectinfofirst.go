package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"net"

	"go.uber.org/zap"
)

var handshake int = 0

func GiveDetectInfoFirst(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// front process back netowork
	logger.Info("GiveDetectInfoFirst: ", zap.Any("message", p.GetMessage()))
	rt := query.First_detect_info(p.GetRkey(), p.GetMessage())
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.UPDATE_DETECT_MODE,
		Message:    rt,
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDetectInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectInfo: ", zap.Any("message", p.GetMessage()))
	if handshake == 0 {
		taskservice.Start()
		handshake = 1
	}
	return task.SUCCESS, nil
}

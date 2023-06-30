package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"net"

	"go.uber.org/zap"
)

func GiveDetectInfoFirst(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	// front process back netowork
	logger.Info("GiveDetectInfoFirst: ", zap.Any("Key", p.GetRkey()), zap.Any("message", p.GetMessage()))
	rt := query.First_detect_info(*Key, p.GetMessage())
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Rkey:       p.GetRkey(),
		Work:       task.UPDATE_DETECT_MODE,
		Message:    rt,
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDetectInfo(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectInfo: ", zap.Any("message", p.GetMessage()))

	return task.SUCCESS, nil
}

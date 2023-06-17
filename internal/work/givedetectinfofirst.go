package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func GiveDetectInfoFirst(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	// front process back netowork
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.UPDATE_DETECT_MODE,
		Message:    "1|1",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	// if p.GetMessage() == "0|0" {
	// 	var send_packet = packet.WorkPacket{
	// 		MacAddress: p.GetMacAddress(),
	// 		IpAddress:  p.GetipAddress(),
	// 		Work:       task.UPDATE_DETECT_MODE,
	// 		Message:    "1|1",
	// 	}
	// 	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	// 	if err != nil {
	// 		return task.FAIL, err
	// 	}
	// 	return task.SUCCESS, nil
	// }
	return task.SUCCESS, nil
}

func GiveDetectInfo(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectInfo: ", zap.Any("message", p.GetMessage()))
	return task.SUCCESS, nil
}

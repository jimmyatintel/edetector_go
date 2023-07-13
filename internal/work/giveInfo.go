package work

import (
	client "edetector_go/internal/client"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"net"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GiveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) { // the first packet: insert user info
	logger.Info("GiveInfo: ", zap.Any("message", p.GetMessage()))
	np := packet.CheckIsWork(p)
	ClientInfo := client.PacketClientInfo(np)
	if (ClientInfo.KeyNum == "") || (ClientInfo.KeyNum == "null") || (ClientInfo.KeyNum == "NoKey") { // assign a new key(uuid)
		ClientInfo.KeyNum = strings.Replace(uuid.New().String(), "-", "", -1)
	}
	query.Checkindex(ClientInfo.KeyNum, p.GetipAddress(), p.GetMacAddress())
	query.Addmachine(ClientInfo)
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.OPEN_CHECK_THREAD,
		Message:    ClientInfo.KeyNum,
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn) // send ack(OPEN_CHECK_THREAD) to client
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDetectPortInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectPortInfo: ", zap.Any("message", p.GetMessage()))
	return task.SUCCESS, nil
}
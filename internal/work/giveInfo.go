package work

import (
	client "edetector_go/internal/client"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	mq "edetector_go/pkg/mariadb/query"
	"net"
	"strings"

	"github.com/google/uuid"
)

func GiveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) { // the first packet: insert user info
	logger.Info("GiveInfo: " + p.GetRkey() + "::" + p.GetMessage())
	logger.Info("ip: " + p.GetipAddress() + ", mac: " + p.GetMacAddress())
	np := packet.CheckIsWork(p)
	ClientInfo, err := client.PacketClientInfo(np)
	if err != nil {
		logger.Error("GiveInfo error: " + err.Error())
		clientsearchsend.SendTCPtoClient(p, task.REJECT_AGENT, "", conn)
		return task.FAIL, err
	}
	if (ClientInfo.KeyNum == "") || (ClientInfo.KeyNum == "null") || (ClientInfo.KeyNum == "NoKey") { // assign a new key(uuid)
		ClientInfo.KeyNum = strings.Replace(uuid.New().String(), "-", "", -1)
		logger.Info("New key: " + ClientInfo.KeyNum)
	}
	mq.Checkindex(ClientInfo.KeyNum, p.GetipAddress(), p.GetMacAddress())
	mq.Addmachine(ClientInfo)
	err = clientsearchsend.SendTCPtoClient(p, task.OPEN_CHECK_THREAD, ClientInfo.KeyNum, conn) // send ack(OPEN_CHECK_THREAD) to client
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

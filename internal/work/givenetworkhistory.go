package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"encoding/json"

	"go.uber.org/zap"
)

type NetworkJson struct {
	PID               int    `json:"pid"`
	Address           string `json:"address"`
	Timestamp         int    `json:"timestamp"`
	ProcessTime       int    `json:"process_time"`
	ConnectionINorOUT bool   `json:"connection_inorout"`
	AgentPort         int    `json:"agent_port"`
}

func (n NetworkJson) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

func GiveNetworkHistory(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveNetworkHistory: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_NETWORK_HISTORY_INFO,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveNetworkHistoryData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveNetworkHistoryData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	//todo
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveNetworkHistoryEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveNetworkHistoryEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))

	// Data := ChangeNetworkToJson(p)
	// template := elasticquery.New_source(p.GetRkey(), "Networkdata")
	// elasticquery.Send_to_elastic("ed_network_history", template, Data)

	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"net"
	"strings"

	"encoding/json"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MemoryNetwork struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	PID               int    `json:"pid"`
	Address           string `json:"address"`
	Timestamp         int    `json:"timestamp"`
	ProcessCreateTime int    `json:"processCreateTime"`
	ConnectionINorOUT bool   `json:"connectionInOrOut"`
	AgentPort         int    `json:"agentPort"`
}

func (n MemoryNetwork) Elastical() ([]byte, error) {
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
	go NetworkElastic(p)
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveNetworkHistoryEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveNetworkHistoryEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	go NetworkElastic(p)
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

func NetworkElastic(p packet.Packet) {
	networkSet := make(map[string]struct{})
	lines := strings.Split(p.GetMessage(), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		line = strings.ReplaceAll(line, "|", "@|@")
		uuid := uuid.NewString()
		values := strings.Split(line, "@|@")
		key := values[0] + "," + values[2]
		networkSet[key] = struct{}{}
		elasticquery.SendToDetailsElastic(uuid, "ed_de_memory_network", p.GetRkey(), line, &MemoryNetwork{}, "ed_high")
	}
	elasticquery.UpdateNetworkInfo(p.GetRkey(), networkSet)
}

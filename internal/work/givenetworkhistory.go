package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"net"

	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

var Data_acache []byte

type NetworkJson struct {
	PID               string `json:"pid"`
	Address           string `json:"address"`
	Timestamp         string `json:"timestamp"`
	ProcessTime       string `json:"process_time"`
	ConnectionINorOUT string `json:"connection_INorOUT"`
	AgentPort         string `json:"agent_port"`
}

func (n NetworkJson) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

func GiveNetworkHistory(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveNetworkHistory: ", zap.Any("message", p.GetMessage()))
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

func GiveNetworkHistoryData(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveNetworkHistoryData: ", zap.Any("message", p.GetMessage()))
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

func GiveNetworkHistoryEnd(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveNetworkHistoryEnd: ", zap.Any("message", p.GetMessage()))
	Data := change2json(p)
	template := elasticquery.New_source(p.GetRkey(), "Networkdata")
	elasticquery.Send_to_elastic(template, Data)
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

func change2json(p packet.Packet) []elasticquery.Request_data {
	lines := strings.Split(p.GetMessage(), "\n")
	var dataSlice []elasticquery.Request_data
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) == 6 {
			data := NetworkJson{
				PID:               values[0],
				Address:           values[1],
				Timestamp:         values[2],
				ProcessTime:       values[3],
				ConnectionINorOUT: values[4],
				AgentPort:         values[5],
			}

			dataSlice = append(dataSlice, elasticquery.Request_data(data))
		}
	}
	jsonData, err := json.Marshal(dataSlice)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return nil
	}
	logger.Debug("Json format: ", zap.Any("json", string(jsonData)))
	return dataSlice
}

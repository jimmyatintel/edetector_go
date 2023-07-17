package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	// "strconv"

	"go.uber.org/zap"
	// "encoding/json"
	// "fmt"
	// "strings"
)

type ProcessDetectJson struct {
	PID         int    `json:"pid"`
	Parent_PID  int    `json:"parent_pid"`
	ProcessName string `json:"process_name"`
	ProcessTime int    `json:"process_time"`
	ParentName  string `json:"parent_name"`
	ParentTime  int    `json:"parent_time"`
}

func (n ProcessDetectJson) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

func GiveProcessHistory(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveProcessHistory: ", zap.Any("message", p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_PROCESS_HISTORY_INFO,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveProcessHistoryData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessHistoryData: ", zap.Any("message", p.GetMessage()))
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

func GiveProcessHistoryEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessHistoryEnd: ", zap.Any("message", p.GetMessage()))
	Data := ChangeProcessToJson(p)
	template := elasticquery.New_source(p.GetRkey(), "Processdata")
	elasticquery.Send_to_elastic("ed_process_history", template, Data)
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

func ChangeProcessToJson(p packet.Packet) []elasticquery.Request_data {
	lines := strings.Split(p.GetMessage(), "\n")
	var dataSlice []elasticquery.Request_data
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		data := ProcessDetectJson{}
		To_json(line, &data)
		dataSlice = append(dataSlice, elasticquery.Request_data(data))
	}
	jsonData, err := json.Marshal(dataSlice)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return nil
	}
	logger.Debug("Json format: ", zap.Any("json", string(jsonData)))
	return dataSlice
}

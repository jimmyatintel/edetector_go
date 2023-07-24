package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"encoding/json"
	"net"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProcessDetectJson struct {
	UUID        string `json:"uuid"`
	AgentID     string `json:"agent_id"`
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
	logger.Info("GiveProcessHistory: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
	logger.Debug("GiveProcessHistoryData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
	logger.Debug("GiveProcessHistoryEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	lines := strings.Split(p.GetMessage(), "\n")
	for _, line := range lines {
		// main index
		values := strings.Split(line, "|")
		elasticID := uuid.NewString()
		template := elasticquery.New_main(elasticID, "network_history", values[1], values[2], "network_record", values[5]) //! ask frontend
		elasticquery.Send_to_main_elastic("main", template)
		// table index
		// data := ProcessDetectJson{}
		// To_json(line, &data)
		// elasticquery.Send_to_elastic("network_history", data)
		// Data := RawDataToJson(p.GetMessage(), ProcessDetectJson{})
		// template := elasticquery.New_source(p.GetRkey(), "Processdata")
		// elasticquery.Send_to_elastic("ed_process_history", template, Data)
	}

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

// func ChangeProcessToJson(p packet.Packet) []elasticquery.Request_data {
// 	lines := strings.Split(p.GetMessage(), "\n")
// 	var dataSlice []elasticquery.Request_data
// 	for _, line := range lines {
// 		if len(line) == 0 {
// 			continue
// 		}
// 		data := ProcessDetectJson{}
// 		To_json(line, &data)
// 		dataSlice = append(dataSlice, elasticquery.Request_data(data))
// 	}
// 	logger.Info("ChangeProcessToJson", zap.Any("message", dataSlice))
// 	return dataSlice
// }

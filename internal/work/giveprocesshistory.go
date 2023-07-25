package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"encoding/json"
	"net"

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

	// send to elasticsearch
	// lines := strings.Split(p.GetMessage(), "\n")
	// for _, line := range lines {
	// 	values := strings.Split(line, "|")
	// 	if len(values) != 6 {
	// 		continue
	// 	}
	// 	uuid := uuid.NewString()
	// 	elasticquery.SendToMainElastic(uuid, "ed_main", p.GetRkey(), values[1], values[2], "network_record", values[5]) //! ask frontend
	// 	// elasticquery.SendToDetailsElastic(uuid, "ed_memory", p.GetRkey(), line, &ProcessDetectJson{})
	// }

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

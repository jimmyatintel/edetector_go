package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
	"encoding/json"
	"fmt"
	"strings"
)

type ProcessJson struct {
	PID                 string `json:"pid"`
	Parent_PID          string `json:"parent_pid"`
	ProcessName         string `json:"process_name"`
	ProcessTime         string `json:"process_time"`
	ParentName          string `json:"parent_name"`
	ParentTime          string `json:"parent_time"`
}

func GiveProcessHistory(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessHistory: ", zap.Any("message", p.GetMessage()))
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

func GiveProcessHistoryData(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
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

func GiveProcessHistoryEnd(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessHistoryEnd: ", zap.Any("message", p.GetMessage()))
	ChangeProcess2Json(p)
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

func ChangeProcess2Json(p packet.Packet) {
	lines := strings.Split(p.GetMessage(), "\n")
	var dataSlice []ProcessJson
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) == 6 {
			data := ProcessJson{
				PID:                values[0],
				Parent_PID:         values[1],
				ProcessName:        values[2],
				ProcessTime:        values[3],
				ParentName:         values[4],
				ParentTime:         values[5],
			}

			dataSlice = append(dataSlice, data)
		}
	}
	jsonData, err := json.Marshal(dataSlice)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}
	logger.Debug("Json format: ", zap.Any("json", string(jsonData)))
}
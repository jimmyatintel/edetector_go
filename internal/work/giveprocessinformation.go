package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
	"encoding/json"
	"fmt"
	"strings"
)

type ProcessInfoJson struct {
	PID            string `json:"pid"`
	ProcessTime    string `json:"process_time"`
	Path           string `json:"path"`
	CommandLine    string `json:"command_line"`
}

func GiveProcessInformation(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveProcessinformation: ", zap.Any("message", p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_PROCESS_INFORMATION,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveProcessInfoData(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessInfoData: ", zap.Any("message", p.GetMessage()))
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

func GiveProcessInfoEnd(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessInfoEnd: ", zap.Any("message", p.GetMessage()))
	ChangeProcessInfo2Json(p)
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

func ChangeProcessInfo2Json(p packet.Packet) {
	lines := strings.Split(p.GetMessage(), "\n")
	var dataSlice []ProcessInfoJson
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) == 4 {
			data := ProcessInfoJson{
				PID:             values[0],
				ProcessTime:     values[1],
				Path:            values[2],
				CommandLine:     values[3],
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
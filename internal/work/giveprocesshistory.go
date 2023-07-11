package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	elasticquery "edetector_go/pkg/elastic/query"
	"net"
	// "strconv"

	"go.uber.org/zap"
	// "encoding/json"
	// "fmt"
	// "strings"
)

type ProcessJson struct {
	PID                 int `json:"pid"`
	Parent_PID          int `json:"parent_pid"`
	ProcessName         string `json:"process_name"`
	ProcessTime         int `json:"process_time"`
	ParentName          string `json:"parent_name"`
	ParentTime          int `json:"parent_time"`
}

func GiveProcessHistory(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
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

func ChangeProcess2Json(p packet.Packet) []elasticquery.Request_data{
	// lines := strings.Split(p.GetMessage(), "\n")
	// var dataSlice []elasticquery.Request_data
	// for _, line := range lines {
	// 	values := strings.Split(line, "|")
	// 	pid, err := strconv.Atoi(values[0])
	// 	if err != nil {
	// 		fmt.Println("Error converting pid to int:", err)
	// 		return nil
	// 	}
	// 	parent_pid, err := strconv.Atoi(values[1])
	// 	if err != nil {
	// 		fmt.Println("Error converting parent_pid to int:", err)
	// 		return nil
	// 	}
	// 	process_time, err := strconv.Atoi(values[3])
	// 	if err != nil {
	// 		fmt.Println("Error converting process_time to int:", err)
	// 		return nil
	// 	}
	// 	parent_time, err := strconv.Atoi(values[5])
	// 	if err != nil {
	// 		fmt.Println("Error converting parent_time to int:", err)
	// 		return nil
	// 	}
	// 	if len(values) == 6 {
	// 		data := ProcessJson{
	// 			PID:                pid,
	// 			Parent_PID:         parent_pid,
	// 			ProcessName:        values[2],
	// 			ProcessTime:        process_time,
	// 			ParentName:         values[4],
	// 			ParentTime:         parent_time,
	// 		}
	// 		dataSlice = append(dataSlice, elasticquery.Request_data(data))
	// 	}
	// }
	// jsonData, err := json.Marshal(dataSlice)
	// if err != nil {
	// 	fmt.Println("Error converting to JSON:", err)
	// 	return nil
	// }
	// logger.Debug("Json format: ", zap.Any("json", string(jsonData)))
	// return dataSlice
	return nil
}
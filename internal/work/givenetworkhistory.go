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

var Data_acache []byte

type NetworkJson struct {
	PID                string `json:"pid"`
	Address            string `json:"address"`
	Timestamp          string `json:"timestamp"`
	ProcessTime        string `json:"process_time"`
	ConnectionINorOUT  string `json:"connection_INorOUT"`
	AgentPort          string `json:"agent_port"`
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

//	func handledata(conn net.Conn) (task.TaskResult, error) {
//		buf := make([]byte, 65536)
//		for {
//			reqLen, err := conn.Read(buf)
//			if err != nil {
//				if err.Error() == "EOF" {
//					logger.Debug("Connection close")
//					return task.FAIL, err
//				} else {
//					logger.Error("Error reading:", zap.Any("error", err.Error()))
//					return task.FAIL, err
//				}
//			}
//			if reqLen > 0 {
//				var NewPacket = new(packet.DataPacket)
//				err := NewPacket.NewPacket(buf)
//				if err != nil {
//					Data_acache = append(Data_acache, buf[:reqLen]...)
//				} else {
//					_, err = WrokMap[NewPacket.GetTaskType()](NewPacket, conn)
//					if err != nil {
//						logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()))
//						return task.FAIL, err
//					}
//					return task.SUCCESS, nil
//				}
//			}
//		}
//	}
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
	change2json(p)
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

func change2json(p packet.Packet) {
	lines := strings.Split(p.GetMessage(), "\n")
	var dataSlice []NetworkJson
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) == 6 {
			data := NetworkJson{
				PID:                values[0],
				Address:            values[1],
				Timestamp:          values[2],
				ProcessTime:        values[3],
				ConnectionINorOUT:  values[4],
				AgentPort:          values[5],
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




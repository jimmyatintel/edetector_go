package workfromapi

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	// "edetector_go/pkg/logger"
	"net"

	// "strings"
	// "go.uber.org/zap"
	// "encoding/json"
	// "strings"
)

type reponseJson struct {
	IsSuccess          string `json:"isSuccess"`
	Message            string `json:"message"`
}

func ChangeDetectMode(p packet.UserPacket, Key *string, conn net.Conn) (task.TaskResult, error) {
	// logger.Info("ChangeDetectMode: ", zap.Any("message", p.GetMessage()))

	// values := strings.Split(p.GetMessage(), "|")
	// msg := reponseJson{
	// 	IsSuccess: values[0],
	// 	Message: values[1],
	// }
	// jsonMsg, err := json.Marshal(msg)

	// var send_packet = packet.TaskPacket{
	// 	MacAddress: p.GetMacAddress(), 
	// 	IpAddress:  p.GetipAddress(),
	// 	Work:       task.GET_NETWORK_HISTORY_INFO,
	// 	Message:    string(jsonMsg),
	// }
	// err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	// if err != nil {
	// 	return task.FAIL, err
	// }
	return task.SUCCESS, nil
}

package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"encoding/json"
	"net"

	// "go.uber.org/zap"
)

func ChangeDetectMode(p packet.UserPacket, Key *string, conn net.Conn) (task.TaskResult, error) {
	// // "0|0"
	// logger.Info("ChangeDetectMode: ", zap.Any("message", p.GetMessage()))

	// "0|0"
	logger.Info("ChangeDetectMode: ", zap.Any("message", p.GetMessage()))

	// "0|0"
	logger.Info("ChangeDetectMode: ", zap.Any("message", p.GetMessage()))

	// Receive request
	requestJSON := make([]byte, 4096)
	n, err := conn.Read(requestJSON)
	if err != nil {
		return task.FAIL, err
	}

	// Parse the JSON request
	var request Request
	err = json.Unmarshal(requestJSON[:n], &request)
	if err != nil {
		return task.FAIL, err
	}

	// Inform agent: "0|0"
	err_agent := clientsearchsend.SendUserTCPtoClient(p.GetRkey())
	if err_agent != nil {
		return task.FAIL, err_agent
	}
	// Generate and send response
	response := Response{
		IsSuccess: true,
		Message:   "msg",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return task.FAIL, err
	}
	_, err = conn.Write(responseJSON)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

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
  
	// "0|0"
	logger.Info("ChangeDetectMode: ", zap.Any("message", p.GetMessage()))
	


	return task.SUCCESS, nil
}

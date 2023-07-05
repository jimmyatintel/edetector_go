package work

import (
	// clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	// elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"net"

	// "encoding/json"
	// "fmt"
	// "strings"

	"go.uber.org/zap"
)


func GiveDriveInfo(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDriveInfo: ", zap.Any("message", p.GetMessage()))
	return task.SUCCESS, nil
}

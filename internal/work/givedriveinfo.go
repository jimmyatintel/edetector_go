package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	// elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"net"

	// "encoding/json"
	// "fmt"
	"strings"

	"go.uber.org/zap"
)


func GiveDriveInfo(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDriveInfo: ", zap.Any("message", p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}

	drives := strings.Split(p.GetMessage(), "|")
	for _, d := range drives {
		parts := strings.Split(d, "-")
		if len(parts) == 2{
			drive := parts[0]
			driveInfo := strings.Split(parts[1], ",")[0]
			msg := drive + "|" + driveInfo
			logger.Info("ExplorerInfo: ", zap.Any("message", msg))
			err = clientsearchsend.SendDriveTCPtoClient(p, p.GetRkey(), task.EXPLORER_INFO, msg + "|Explorer|ScheduleName|0|2048")
			if err != nil {
				return task.FAIL, err
			}
		}
	}
	return task.SUCCESS, nil
}
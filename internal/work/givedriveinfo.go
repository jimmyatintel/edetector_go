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
	// "strings"

	"go.uber.org/zap"
)


func GiveDriveInfo(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDriveInfo: ", zap.Any("message", p.GetMessage()))
	// parts := strings.Split(p.GetMessage(), "-")
	// drive := parts[0]
	// driveInfo := strings.Split(parts[1], ",")[0]

	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.EXPLORER_INFO,
		Message:    "null",
		// Message:    drive + "|" + driveInfo + "|Explorer|ScheduleName",
	}
	// fmt.Println(string(task.TRANSPORT_EXPLORER))
	// fmt.Println(drive + "|" + driveInfo + "|Explorer|ScheduleName")
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}
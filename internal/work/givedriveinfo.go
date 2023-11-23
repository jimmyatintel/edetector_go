package work

import (
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"

	"net"
)

func GiveDriveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDriveInfo: " + p.GetRkey() + "::" + p.GetMessage())
	elastic.DeleteByQueryRequest("agent", p.GetRkey(), "StartGetDrive")
	go HandleExpolorer(p)
	return task.SUCCESS, nil
}

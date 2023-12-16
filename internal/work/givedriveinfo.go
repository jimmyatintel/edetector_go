package work

import (
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"fmt"

	"net"
)

func GiveDriveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDriveInfo: " + p.GetRkey() + "::" + p.GetMessage())
	deleteQuery := fmt.Sprintf(`
	{
		"query": {
			"term": {
				"agent": "%s"
			}
		}
	}
	`, p.GetRkey())
	elastic.DeleteByQueryRequest(deleteQuery, "StartGetDrive")
	go HandleExpolorer(p)
	return task.SUCCESS, nil
}

package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"

	"net"
)

func GiveDriveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDriveInfo: " + p.GetRkey() + "||" + p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	elastic.DeleteByQueryRequest("agent", p.GetRkey(), "StartGetDrive")
	go HandleExpolorer(p)
	return task.SUCCESS, nil
}

package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
)

func StartScan(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartScan: " + p.GetRkey() + "|" + p.GetMessage())
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_SCAN, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

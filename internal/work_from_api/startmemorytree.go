package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
)

func StartMemoryTree(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartMemoryTree: " + p.GetRkey())

	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_MEMORY_TREE, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

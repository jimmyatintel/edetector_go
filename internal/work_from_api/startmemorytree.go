package workfromapi

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
)

func StartMemoryTree(p packet.UserPacket) (task.TaskResult, error) {
	// logger.Info("StartMemoryTree: " + p.GetRkey())
	// err := clientsearchsend.SendUserTCPtoClient(p, task.GET_MEMORY_TREE, p.GetMessage())
	// if err != nil {
	// 	return task.FAIL, err
	// }
	return task.SUCCESS, nil
}

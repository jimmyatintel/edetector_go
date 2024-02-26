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

func StartLoadDll(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartLoadDll: " + p.GetRkey())
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_LOAD_DLL, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func StartDumpDll(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartDumpDll: " + p.GetRkey())
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_DUMP_DLL, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func StartDumpProcess(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartDumpProcess: " + p.GetRkey())
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_DUMP_PROCESS, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}
package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
)

func StartLoadDll(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info(p.GetRkey() + "::StartLoadDll: " + p.GetMessage())
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_LOAD_DLL, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func StartDumpDll(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info(p.GetRkey() + "::StartDumpDll: " + p.GetMessage())
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_DUMP_DLL, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func StartDumpProcess(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info(p.GetRkey() + "::StartDumpProcess: " + p.GetMessage())
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_DUMP_PROCESS, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func StartDumpDrive(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info(p.GetRkey() + "::StartDumpDrive: " + p.GetMessage())
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_DUMP_DRIVE, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

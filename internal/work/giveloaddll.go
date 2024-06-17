package work

import (
	"edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"net"
)

func GiveLoadDllData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	pathInfo := p.GetMessage()
	logger.Debug(key + "::GiveLoadDllData: " + pathInfo)

	// send path info in load task channel to trigger response
	load_chan, err := channelmap.GetLoadDumpChannel(key + string(task.START_LOAD_DLL) + p.GetMessage())
	if err != nil {
		logger.Error("Error getting load channel: " + err.Error())
		return task.FAIL, err
	}
	load_chan <- pathInfo

	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveLoadDllEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveLoadDllEnd: " + key)
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	query.Finish_task(key, "StartLoadDll")
	return task.SUCCESS, nil
}

package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/connectionmap"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"edetector_go/pkg/request"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

func GiveDumpProcessInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	msg, msgs := p.GetMessage(), strings.Split(p.GetMessage(), "|")
	logger.Info(key + "::GiveDumpProcessInfo: " + msg)

	dataLen, pid := msgs[0], msgs[1]
	if dataLen == "-1" {
		logger.Error(key + "::GiveDumpProcessInfo: Dump process path not found")
	}

	total, err := strconv.Atoi(dataLen)
	if err != nil {
		return task.FAIL, err
	}

	// get taskId from redis
	taskId, err := redis.RedisGetString(key + string(task.START_DUMP_PROCESS) + pid)
	if err != nil {
		return task.FAIL, err
	}

	connectionmap.StoreConnInfo(conn, connectionmap.ConnInfo{
		TaskId:  taskId,
		Msg:     pid,
		DataLen: total,
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveDumpProcessData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpProcessData")

	// get connInfo from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	// write file
	path := filepath.Join(dumpWorkingPath, connInfo.TaskId)
	content := getDataPacketContent(p)
	if err := file.WriteFile(path, content); err != nil {
		return task.FAIL, err
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveDumpProcessEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpProcessEnd")

	// remove the msg from ConnMsgMap before return
	defer connectionmap.RemoveConnInfo(conn)

	// get connInfo from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	if connInfo.DataLen != -1 {
		srcPath := filepath.Join(dumpWorkingPath, connInfo.TaskId)
		destPath := filepath.Join(dumpUstagePath, connInfo.TaskId+".zip")

		// truncate data
		if err := file.TruncateFile(srcPath, connInfo.DataLen); err != nil {
			logger.Error("TruncateFile: " + err.Error())
			return task.FAIL, err
		}

		// move to unstage
		if err := file.MoveFile(srcPath, destPath); err != nil {
			logger.Error("MoveFile: " + err.Error())
			return task.FAIL, err
		}

		// inform API that dump file is ready
		request.LoadDumpReady(request.ReadyData{
			TaskId: connInfo.TaskId,
		})
	} else {
		// inform API that there is an error
		request.LoadDumpReady(request.ReadyData{
			TaskId: connInfo.TaskId,
			Failed: "Dump process pid not found",
		})
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("SendTCPtoClient: " + err.Error())
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

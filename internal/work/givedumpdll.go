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

func GiveDumpDllInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	msg, msgs := p.GetMessage(), strings.Split(p.GetMessage(), "|")
	logger.Info(key + "::GiveDumpDllInfo: " + msg)

	dataLenStr, pid, path := msgs[0], msgs[1], msgs[2]
	if dataLenStr == "-1" {
		logger.Error(key + "::GiveDumpDllInfo: Dump dll path not found")
	}

	dataLen, err := strconv.Atoi(dataLenStr)
	if err != nil {
		return task.FAIL, err
	}

	// get taskId from redis
	taskId, err := redis.RedisGetString(key + string(task.START_DUMP_DLL) + pid + path)
	if err != nil {
		logger.Error("Error getting taskId from redis: " + err.Error())
		return task.FAIL, err
	}

	connectionmap.StoreConnInfo(conn, connectionmap.ConnInfo{
		TaskId:  taskId,
		Msg:     pid + path,
		DataLen: dataLen,
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveDumpDllData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpDllData")

	// get dll path from ConnMsgMap
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

func GiveDumpDllEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpDllEnd")

	// remove the msg from ConnMsgMap before return
	defer connectionmap.RemoveConnInfo(conn)

	// get dll path from ConnMsgMap
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
			logger.Error("Error truncating file: " + err.Error())
			return task.FAIL, err
		}

		// move to unstage
		if err := file.MoveFile(srcPath, destPath); err != nil {
			logger.Error("Error moving file: " + err.Error())
			return task.FAIL, err
		}

		// inform API that the dump dll is ready
		request.LoadDumpReady(request.ReadyData{
			TaskId: connInfo.TaskId,
		})
	} else {
		// inform API that there is an error
		request.LoadDumpReady(request.ReadyData{
			TaskId: connInfo.TaskId,
			Failed: "Dump dll path and pid not found",
		})
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("Error sending packet: " + err.Error())
		return task.FAIL, err
	}

	// remove the taskId from redis
	if err := redis.RedisDelete(key + string(task.START_DUMP_DLL) + connInfo.Msg); err != nil {
		logger.Error("Error deleting key from redis: " + err.Error())
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

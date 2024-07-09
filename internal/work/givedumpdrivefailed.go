package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/connectionmap"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"edetector_go/pkg/request"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GiveDumpDriveFailedInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key, msg := p.GetRkey(), p.GetMessage()
	logger.Info(key + "::GiveDumpDriveFailedInfo: " + msg)

	dataLen, err := strconv.Atoi(strings.Split(msg, "|")[0])
	if err != nil {
		logger.Error("Error converting dataLen to int: " + err.Error())
		return task.FAIL, err
	}

	// update data length in ConnMsgMap
	oldInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	connectionmap.StoreConnInfo(conn, connectionmap.ConnInfo{
		TaskId:  oldInfo.TaskId,
		Msg:     oldInfo.Msg,
		DataLen: dataLen,
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveDumpDriveFailedData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpDriveFailedData")

	// get taskId from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	// write file
	path := filepath.Join(dumpWorkingPath, connInfo.TaskId+"-failed.zip")
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

func GiveDumpDriveFailedEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpDriveFailedEnd")

	// remove the msg from ConnMsgMap before return
	defer connectionmap.RemoveConnInfo(conn)

	// get taskId from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("SendTCPtoClient: " + err.Error())
		return task.FAIL, err
	}

	// decompress file
	srcPath := filepath.Join(dumpWorkingPath, connInfo.TaskId+"-failed.zip")
	destPath := filepath.Join(dumpWorkingPath, connInfo.TaskId+"-failed.txt")
	if err := file.DecompressFile(srcPath, destPath, connInfo.DataLen); err != nil {
		return task.FAIL, err
	}

	// read error path and send to API
	failed := "Paths failed: "
	if content, err := file.ReadFileLineByLine(destPath); err != nil {
		return task.FAIL, err
	} else {
		failed += strings.Join(content, ",")
	}
	request.LoadDumpReady(request.ReadyData{
		TaskId: connInfo.TaskId,
		Failed: failed,
	})

	// remove the file
	if err := os.Remove(destPath); err != nil {
		logger.Error("Error removing file: " + err.Error())
	}

	// remove key from redis
	if err := redis.RedisDelete(key + string(task.START_DUMP_DRIVE) + connInfo.Msg); err != nil {
		logger.Error("Error deleting key from redis: " + err.Error())
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

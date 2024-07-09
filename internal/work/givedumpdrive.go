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

func GiveDumpDriveProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key, msg := p.GetRkey(), p.GetMessage()
	logger.Info(key + "::GiveDumpDriveProgress: " + msg)

	// update progress
	msgs := strings.Split(msg, "/")
	finished, err := strconv.Atoi(msgs[0])
	if err != nil {
		logger.Error("Error converting finished to int: " + err.Error())
		return task.FAIL, err
	}
	total, err := strconv.Atoi(msgs[1])
	if err != nil {
		logger.Error("Error converting total to int: " + err.Error())
		return task.FAIL, err
	}

	// update progress set in redis and inform API (POST to loadDumpReady) that progress updated
	progress := float64(finished) / float64(total) * 100
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}
	redis.UpdateDumpProgress(connInfo.TaskId, int(progress))
	request.LoadDumpReady(request.ReadyData{TaskId: connInfo.TaskId})

	// send data right msg to client
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		logger.Error("SendTCPtoClient: " + err.Error())
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveDumpDriveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpDriveInfo: " + p.GetMessage())

	dataLen, err := strconv.Atoi(strings.Split(p.GetMessage(), "|")[0])
	if err != nil {
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

func GiveDumpDriveData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpDriveData")

	// get taskId from ConnMsgMap
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

func GiveDumpDriveEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpDriveEnd")

	// get taskId from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	workPath := filepath.Join(dumpWorkingPath, connInfo.TaskId)
	unstagePath := filepath.Join(dumpUstagePath, connInfo.TaskId+".zip")

	// truncate data
	if err := file.TruncateFile(workPath, connInfo.DataLen); err != nil {
		logger.Error("TruncateFile: " + err.Error())
		return task.FAIL, err
	}

	// move to unstage
	if err := file.MoveFile(workPath, unstagePath); err != nil {
		logger.Error("MoveFile: " + err.Error())
		return task.FAIL, err
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("SendTCPtoClient: " + err.Error())
		return task.FAIL, err
	}

	// update progress set in redis and inform API (POST to loadDumpReady) that progress updated
	redis.UpdateDumpProgress(connInfo.TaskId, 100)
	request.LoadDumpReady(request.ReadyData{TaskId: connInfo.TaskId})

	return task.SUCCESS, nil
}

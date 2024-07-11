package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/connectionmap"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
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

	// get taskId from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	dataLen, err := strconv.Atoi(strings.Split(msg, "|")[0])
	if err != nil {
		logger.Error("Error converting dataLen to int: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	// update data length in ConnMsgMap
	connectionmap.StoreConnInfo(conn, connectionmap.ConnInfo{
		TaskId:  connInfo.TaskId,
		Msg:     connInfo.Msg,
		DataLen: dataLen,
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("SendTCPtoClient: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
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
		logger.Error("Error writing file: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("SendTCPtoClient: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
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
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	// decompress file
	srcPath := filepath.Join(dumpWorkingPath, connInfo.TaskId+"-failed.zip")
	destPath := filepath.Join(dumpWorkingPath, connInfo.TaskId+"-failed.txt")
	if err := file.DecompressFile(srcPath, destPath, connInfo.DataLen); err != nil {
		logger.Error("Error decompressing file: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	// read error path and send to API
	content, err := file.ReadFileLineByLine(destPath)
	if err != nil {
		logger.Error("Error reading file: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}
	request.LoadDumpReady(request.ReadyData{
		TaskId:   connInfo.TaskId,
		Failed:   "Partial paths failed: " + strings.Join(content, ","),
		RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
		Progress: -2,
	})

	// remove the file
	if err := os.Remove(destPath); err != nil {
		logger.Error("Error removing file: " + err.Error())
	}

	return task.SUCCESS, nil
}

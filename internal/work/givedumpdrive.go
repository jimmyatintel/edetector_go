package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/connectionmap"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/request"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

func GiveDumpDriveProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key, msg := p.GetRkey(), p.GetMessage()
	logger.Info(key + "::GiveDumpDriveProgress: " + msg)

	// get connInfo from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	// update progress
	msgs := strings.Split(msg, "/")
	finished, err := strconv.Atoi(msgs[0])
	if err != nil {
		logger.Error("Error converting finished to int: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}
	total, err := strconv.Atoi(msgs[1])
	if err != nil {
		logger.Error("Error converting total to int: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	// update progress set in redis and inform API (POST to loadDumpReady) that progress updated
	progress := float64(finished) / float64(total) * config.Viper.GetFloat64("DUMP_FIRST_PART")
	request.LoadDumpReady(request.ReadyData{
		TaskId:   connInfo.TaskId,
		Progress: int(progress),
	})

	// send data right msg to client
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
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

func GiveDumpDriveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveDumpDriveInfo: " + p.GetMessage())

	// get connInfo from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	dataLen, err := strconv.Atoi(strings.Split(p.GetMessage(), "|")[0])
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
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	// update Progress
	progress := config.Viper.GetFloat64("DUMP_FIRST_PART") + float64(len(content))/float64(connInfo.DataLen)*(100-config.Viper.GetFloat64("DUMP_FIRST_PART"))
	request.LoadDumpReady(request.ReadyData{
		TaskId:   connInfo.TaskId,
		Progress: int(progress),
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
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	// move to unstage
	if err := file.MoveFile(workPath, unstagePath); err != nil {
		logger.Error("MoveFile: " + err.Error())
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

	// update progress set in redis and inform API (POST to loadDumpReady) that progress updated
	request.LoadDumpReady(request.ReadyData{
		TaskId:   connInfo.TaskId,
		Progress: 100,
	})

	return task.SUCCESS, nil
}

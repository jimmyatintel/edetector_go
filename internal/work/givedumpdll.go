package work

import (
	"edetector_go/config"
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

func ReadyDumpDll(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key, msg := p.GetRkey(), p.GetMessage()
	logger.Info(key + "::ReadyDumpDll: " + msg)

	// get taskId from redis
	taskId, err := redis.RedisGetString(key + string(task.START_DUMP_DLL) + msg)
	if err != nil {
		logger.Error("Error getting taskId from redis: " + err.Error())
		return task.FAIL, err
	}

	// store taskId and msg to ConnMsgMap
	connectionmap.StoreConnInfo(conn, connectionmap.ConnInfo{
		TaskId:  taskId,
		Msg:     msg,
		DataLen: 0,
	})

	return task.SUCCESS, nil
}

func GiveDumpDllProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key, msg := p.GetRkey(), p.GetMessage()
	logger.Info(key + "::GiveDumpDllProgress: " + msg)

	// get connInfo from ConnMsg
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	// count progress
	msgs := strings.Split(msg, "/")
	finished, err := strconv.Atoi(msgs[0])
	if err != nil {
		logger.Error("Error converting finished to int: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
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
			Progress: -1,
		})
		return task.FAIL, err
	}

	// update progress
	progress := float64(finished) / float64(total) * config.Viper.GetFloat64("DUMP_FIRST_PART")
	request.LoadDumpReady(request.ReadyData{
		TaskId:   connInfo.TaskId,
		Progress: int(progress),
	})

	// send data right to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("Error sending packet: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveDumpDllInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	msg, msgs := p.GetMessage(), strings.Split(p.GetMessage(), "|")
	logger.Info(key + "::GiveDumpDllInfo: " + msg)

	dataLenStr := msgs[0]
	if dataLenStr == "-1" {
		logger.Error(key + "::GiveDumpDllInfo: Dump dll path not found")
	}

	// get connInfo from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	dataLen, err := strconv.Atoi(dataLenStr)
	if err != nil {
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	connectionmap.StoreConnInfo(conn, connectionmap.ConnInfo{
		TaskId:  connInfo.TaskId,
		Msg:     connInfo.Msg,
		DataLen: dataLen,
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
			Progress: -1,
		})
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
		logger.Error("Error writing file: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
			Progress: -1,
		})
		return task.FAIL, err
	}

	// update progress
	progress := config.Viper.GetFloat64("DUMP_FIRST_PART") + float64(len(content))/float64(connInfo.DataLen)*(100-config.Viper.GetFloat64("DUMP_FIRST_PART"))
	request.LoadDumpReady(request.ReadyData{
		TaskId:   connInfo.TaskId,
		Progress: int(progress),
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("Error sending packet: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
			Progress: -1,
		})
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
			request.LoadDumpReady(request.ReadyData{
				TaskId:   connInfo.TaskId,
				Failed:   "InternalServerError",
				RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
				Progress: -1,
			})
			return task.FAIL, err
		}

		// move to unstage
		if err := file.MoveFile(srcPath, destPath); err != nil {
			logger.Error("Error moving file: " + err.Error())
			request.LoadDumpReady(request.ReadyData{
				TaskId:   connInfo.TaskId,
				Failed:   "InternalServerError",
				RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
				Progress: -1,
			})
			return task.FAIL, err
		}

		// inform API that the dump dll is ready
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
			Progress: 100,
		})
	} else {
		// inform API that there is an error
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "Dump dll path and pid not found",
			RedisKey: key + string(task.START_DUMP_DLL) + connInfo.Msg,
			Progress: -1,
		})
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("Error sending packet: " + err.Error())
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

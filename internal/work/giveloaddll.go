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
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GiveLoadDllInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	msg, msgs := p.GetMessage(), strings.Split(p.GetMessage(), "|")
	dataLen, pid := msgs[0], msgs[1]+"|"+msgs[2]
	logger.Info(key + "::GiveLoadDllInfo: " + msg)

	// get taskId from the redis
	taskId, err := redis.RedisGetString(key + string(task.START_LOAD_DLL) + pid)
	if err != nil {
		logger.Error(key + "::GiveLoadDllInfo: " + err.Error())
		return task.FAIL, err
	}

	if dataLen == "-1" {
		logger.Error(key + "::GiveLoadDllInfo: pid not found")
	} else if dataLen == "0" {
		logger.Warn(key + "::GiveLoadDllInfo: dll not found for pid = " + pid)
	}

	total, err := strconv.Atoi(dataLen)
	if err != nil {
		logger.Error("Error convert dataLen to int" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_LOAD_DLL) + pid,
			LoadDll:  true,
			Progress: -1,
		})
		return task.FAIL, err
	}

	connectionmap.StoreConnInfo(conn, connectionmap.ConnInfo{
		TaskId:     taskId,
		Msg:        pid,
		DataLen:    total,
		CurDataLen: 0,
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("Error sending data right msg to client: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_LOAD_DLL) + pid,
			LoadDll:  true,
			Progress: -1,
		})
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveLoadDllData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveLoadDllData")

	// get pid from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting pid from ConnMsgMap")
		return task.FAIL, nil
	}

	// write file
	path := filepath.Join(loadDllWorkingPath, connInfo.TaskId+".zip")
	content := getDataPacketContent(p)
	if err := file.WriteFile(path, content); err != nil {
		logger.Error("Error writing file: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_LOAD_DLL) + connInfo.Msg,
			LoadDll:  true,
			Progress: -1,
		})
		return task.FAIL, err
	}

	curDataLen := connectionmap.UpdateConnInfoCurDataLen(conn, len(content))
	progress := config.Viper.GetFloat64("DUMP_FIRST_PART") + float64(curDataLen)/float64(connInfo.DataLen)*(99-config.Viper.GetFloat64("DUMP_FIRST_PART"))
	request.LoadDumpReady(request.ReadyData{
		TaskId:   connInfo.TaskId,
		Progress: int(progress),
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("Error sending data right msg to client: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_LOAD_DLL) + connInfo.Msg,
			LoadDll:  true,
			Progress: -1,
		})
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveLoadDllEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveLoadDllEnd")

	// remove the msg from ConnMsgMap before return
	defer connectionmap.RemoveConnInfo(conn)

	// get the pid and path info from ConnMsgMap
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	if connInfo.DataLen > 0 {
		workPath := filepath.Join(loadDllWorkingPath, connInfo.TaskId+".zip")
		unstagePath := filepath.Join(loadDllWorkingPath, connInfo.TaskId)

		// truncate data
		if err := file.TruncateFile(workPath, connInfo.DataLen); err != nil {
			logger.Error("Error truncating file: " + err.Error())
			request.LoadDumpReady(request.ReadyData{
				TaskId:   connInfo.TaskId,
				Failed:   "InternalServerError",
				RedisKey: key + string(task.START_LOAD_DLL) + connInfo.Msg,
				LoadDll:  true,
				Progress: -1,
			})
			return task.FAIL, err
		}

		// decompress the file
		if err := file.DecompressFile(workPath, unstagePath, connInfo.DataLen); err != nil {
			logger.Error("Error unzipping file: " + err.Error())
			request.LoadDumpReady(request.ReadyData{
				TaskId:   connInfo.TaskId,
				Failed:   "InternalServerError",
				RedisKey: key + string(task.START_LOAD_DLL) + connInfo.Msg,
				LoadDll:  true,
				Progress: -1,
			})
			return task.FAIL, err
		}

		// read the file line by line as string array
		lines, err := file.ReadFileLineByLine(unstagePath)
		if err != nil {
			logger.Error("Error reading file: " + err.Error())
			request.LoadDumpReady(request.ReadyData{
				TaskId:   connInfo.TaskId,
				Failed:   "InternalServerError",
				RedisKey: key + string(task.START_LOAD_DLL) + connInfo.Msg,
				LoadDll:  true,
				Progress: -1,
			})
			return task.FAIL, err
		}

		// send path info to API
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			DllPaths: strings.Join(lines, "|"),
			RedisKey: key + string(task.START_LOAD_DLL) + connInfo.Msg,
			LoadDll:  true,
			Progress: 100,
		})

		// remove the working file
		if err := os.Remove(unstagePath); err != nil {
			logger.Error("Error removing file: " + err.Error())
			return task.FAIL, err
		}
	} else {
		// send error to API
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   key + "::GiveLoadDllInfo: pid not found",
			RedisKey: key + string(task.START_LOAD_DLL) + connInfo.Msg,
			LoadDll:  true,
			Progress: -1,
		})
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

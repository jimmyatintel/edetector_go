package work

import (
	"edetector_go/config"
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
	connectionmap.UpdateConnInfoDataLen(conn, dataLen)

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

	curDataLen := connectionmap.UpdateConnInfoCurDataLen(conn, len(content))

	// update Progress
	progress := config.Viper.GetFloat64("DUMP_DRIVE_THIRD_PART") + float64(curDataLen)/float64(connInfo.DataLen)*(config.Viper.GetFloat64("DUMP_DRIVE_FOURTH_PART")-config.Viper.GetFloat64("DUMP_DRIVE_THIRD_PART"))
	request.LoadDumpReady(request.ReadyData{
		TaskId:   connInfo.TaskId,
		Progress: int(progress),
	})
	logger.Info(key + "::GiveDumpDriveFailedData: update progress to " + strconv.Itoa(int(progress)))

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
	dumpPathPath := filepath.Join(dumpWorkingPath, connInfo.TaskId+".txt")
	srcPath := filepath.Join(dumpWorkingPath, connInfo.TaskId+"-failed.zip")
	destPath := filepath.Join(dumpWorkingPath, connInfo.TaskId+"-failed.txt")
	var failedLines int = connInfo.DataLen

	if failedLines > 0 {
		var err error
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

		failedLines, err = file.GetNumberOfLine(destPath)
		if err != nil {
			logger.Error("Error getting number of lines: " + err.Error())
			request.LoadDumpReady(request.ReadyData{
				TaskId:   connInfo.TaskId,
				Failed:   "InternalServerError",
				RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
				Progress: -1,
			})
		}
	}

	pathLines, err := file.GetNumberOfLine(dumpPathPath)
	if err != nil {
		logger.Error("Error getting number of lines: " + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "InternalServerError",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
	}

	// check if all paths failed
	if failedLines == pathLines {
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			Failed:   "All paths failed",
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: -1,
		})
	} else if failedLines > 0 {
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
	} else {
		request.LoadDumpReady(request.ReadyData{
			TaskId:   connInfo.TaskId,
			RedisKey: key + string(task.START_DUMP_DRIVE) + connInfo.Msg,
			Progress: 100,
		})
		logger.Info(key + "::GiveDumpDriveFailedEnd: update progress to " + strconv.Itoa(100))
	}

	// remove the file
	if failedLines > 0 {
		if err := os.Remove(destPath); err != nil {
			logger.Error("Error removing file: " + err.Error())
		}
	}
	if err := os.Remove(dumpPathPath); err != nil {
		logger.Error("Error removing txt in working path: " + err.Error())
	}

	return task.SUCCESS, nil
}

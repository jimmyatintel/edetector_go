package work

import (
	"edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/connectionmap"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
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
	dataLen, pid := msgs[0], msgs[1]
	logger.Info(key + "::GiveLoadDllInfo: " + msg)

	if dataLen == "-1" {
		logger.Error(key + "::GiveLoadDllInfo: pid not found")
	} else if dataLen == "0" {
		logger.Warn(key + "::GiveLoadDllInfo: dll not found for pid = " + pid)
	}

	total, err := strconv.Atoi(dataLen)
	if err != nil {
		return task.FAIL, err
	}

	connectionmap.StoreConnMsg(conn, connectionmap.ConnMsg{
		Info:    pid,
		DataLen: total,
	})

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveLoadDllData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveLoadDllData")

	// get pid from ConnMsgMap
	msg, ok := connectionmap.GetConnMsg(conn)
	if !ok {
		logger.Error("Error getting pid from ConnMsgMap")
		return task.FAIL, nil
	}

	// write file
	path := filepath.Join(loadDllWorkingPath, key+"-"+msg.Info)
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

func GiveLoadDllEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	logger.Info(key + "::GiveLoadDllEnd")

	// get the pid and path info from ConnMsgMap
	workPath, unstagePath, returnData := "", "", ""
	msg, ok := connectionmap.GetConnMsg(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	if msg.DataLen > 0 {
		workPath = filepath.Join(loadDllWorkingPath, key+"-"+msg.Info)
		unstagePath = filepath.Join(loadDllUstagePath, key+"-"+msg.Info)

		// truncate data
		if err := file.TruncateFile(workPath, msg.DataLen); err != nil {
			logger.Error("Error truncating file: " + err.Error())
			return task.FAIL, err
		}

		// decompress the file
		if err := file.DecompressFile(workPath, unstagePath, msg.DataLen); err != nil {
			logger.Error("Error unzipping file: " + err.Error())
			return task.FAIL, err
		}

		// read the file line by line as string array
		lines, err := file.ReadFileLineByLine(unstagePath)
		if err != nil {
			logger.Error("Error reading file: " + err.Error())
			return task.FAIL, err
		}

		returnData = strings.Join(lines, "|")

		// remove the working file
		if err := os.Remove(unstagePath); err != nil {
			logger.Error("Error removing file: " + err.Error())
			return task.FAIL, err
		}
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		return task.FAIL, err
	}

	// send path info in load task channel to trigger response
	load_chan, err := channelmap.GetLoadDumpChannel(key + "-" + string(task.START_LOAD_DLL) + "-" + msg.Info)
	if err != nil {
		logger.Error("Error getting load channel: " + err.Error())
		return task.FAIL, err
	}
	load_chan <- returnData

	// remove the msg from ConnMsgMap
	connectionmap.RemoveConnMsg(conn)

	return task.SUCCESS, nil
}

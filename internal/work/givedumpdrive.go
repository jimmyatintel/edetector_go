package work

import (
	channelmap "edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/connectionmap"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

func GiveDumpDriveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// retrieve data from packet
	key := p.GetRkey()
	msg, msgs := p.GetMessage(), strings.Split(p.GetMessage(), "|")
	logger.Info(key + "::GiveDumpDriveInfo: " + msg)

	dataLen, path := msgs[0], msgs[1]
	if dataLen == "-1" {
		logger.Error(key + "::GiveDumpDriveInfo: Dump drive path not found")
	}

	total, err := strconv.Atoi(dataLen)
	if err != nil {
		return task.FAIL, err
	}

	connectionmap.StoreConnMsg(conn, connectionmap.ConnMsg{
		Info:    path,
		DataLen: total,
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

	// get path from ConnMsgMap
	msg, ok := connectionmap.GetConnMsg(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	// write file
	path := filepath.Join(dumpDriveWorkingPath, key+"-"+msg.Info)
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

	// get path from ConnMsgMap
	workPath, unstagePath := "", ""
	msg, ok := connectionmap.GetConnMsg(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return task.FAIL, nil
	}

	if msg.DataLen != -1 {
		workPath = filepath.Join(dumpDriveWorkingPath, key+"-"+msg.Info)
		unstagePath = filepath.Join(dumpDriveUstagePath, key+"-"+msg.Info+".zip")

		// truncate data
		if err := file.TruncateFile(workPath, msg.DataLen); err != nil {
			logger.Error("TruncateFile: " + err.Error())
			return task.FAIL, err
		}

		// move to unstage
		if err := file.MoveFile(workPath, unstagePath); err != nil {
			logger.Error("MoveFile: " + err.Error())
			return task.FAIL, err
		}
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("SendTCPtoClient: " + err.Error())
		return task.FAIL, err
	}

	// send dump file name in dump task channel to trigger response
	dump_chan, err := channelmap.GetLoadDumpChannel(key + "-" + string(task.START_DUMP_DRIVE) + "-" + msg.Info)
	if err != nil {
		logger.Error("Error getting dump channel: " + err.Error())
		return task.FAIL, err
	}
	dump_chan <- unstagePath

	// remove the msg from ConnMsgMap
	connectionmap.RemoveConnMsg(conn)

	return task.SUCCESS, nil
}

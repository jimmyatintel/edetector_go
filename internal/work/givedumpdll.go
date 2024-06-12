package work

import (
	"edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"net"
	"path/filepath"
	"strconv"
)

func GiveDumpDllInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveDumpDllInfo: " + key + "::" + p.GetMessage())
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-DumpDllTotal", total)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDumpDllData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveDumpDllData: " + key)
	// write file
	path := filepath.Join(dumpDllWorkingPath, key)
	content := getDataPacketContent(p)
	err := file.WriteFile(path, content)
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDumpDllEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveDumpDllEnd: " + key)
	workPath := filepath.Join(dumpDllWorkingPath, key)
	unstagePath := filepath.Join(dumpDllUstagePath, key+".zip")

	// truncate data
	if err := file.TruncateFile(workPath, redis.RedisGetInt(key+"-DumpDllTotal")); err != nil {
		logger.Error("Error truncating file: " + err.Error())
		return task.FAIL, err
	}

	// move to unstage
	if err := file.MoveFile(workPath, unstagePath); err != nil {
		logger.Error("Error moving file: " + err.Error())
		return task.FAIL, err
	}

	// send data right msg to client
	if err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn); err != nil {
		logger.Error("Error sending packet: " + err.Error())
		return task.FAIL, err
	}

	// send dump file name in dump task channel to trigger response
	dump_chan, err := channelmap.GetLoadDumpChannel(key + string(task.START_DUMP_DLL))
	if err != nil {
		logger.Error("Error getting dump channel: " + err.Error())
		return task.FAIL, err
	}
	dump_chan <- unstagePath

	return task.SUCCESS, nil
}

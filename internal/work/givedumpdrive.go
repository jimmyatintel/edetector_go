package work

import (
	channelmap "edetector_go/internal/channelmap"
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

func GiveDumpDriveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveDumpDriveInfo: " + key + "::" + p.GetMessage())
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-DumpDriveTotal", total)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDumpDriveData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveDumpDriveData: " + key)
	// write file
	path := filepath.Join(dumpDriveWorkingPath, key)
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

func GiveDumpDriveEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveDumpDriveEnd: " + key)
	workPath := filepath.Join(dumpDriveWorkingPath, key)
	unstagePath := filepath.Join(dumpDriveUstagePath, key+".zip")

	// truncate data
	if err := file.TruncateFile(workPath, redis.RedisGetInt(key+"-DumpDriveTotal")); err != nil {
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

	// send dump file name in dump task channel to trigger response
	dump_chan, err := channelmap.GetLoadDumpChannel(key + string(task.START_DUMP_DRIVE))
	if err != nil {
		logger.Error("Error getting dump channel: " + err.Error())
		return task.FAIL, err
	}
	dump_chan <- unstagePath

	return task.SUCCESS, nil
}

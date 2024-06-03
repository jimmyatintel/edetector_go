package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
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
	err := file.TruncateFile(workPath, redis.RedisGetInt(key+"-DumpDriveTotal"))
	if err != nil {
		return task.FAIL, err
	}
	// move to unstage
	err = file.MoveFile(workPath, unstagePath)
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	query.Finish_task(key, "StartDumpDrive")
	return task.SUCCESS, nil
}

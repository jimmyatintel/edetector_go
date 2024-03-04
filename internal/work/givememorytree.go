package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

var memoryTreeWorkingPath = "memoryTreeWorking"
var memoryTreeUstagePath = "memoryTreeUnstage"

func init() {
	file.ClearDirContent(memoryTreeWorkingPath)
	file.CheckDir(memoryTreeUstagePath)
}

func GiveMemoryTreeInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveMemoryTreeInfo: " + key + "::" + p.GetMessage())
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-MemoryTreeTotal", total)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveMemoryTree(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveMemoryTree: " + key)
	// write file
	path := filepath.Join(memoryTreeWorkingPath, key)
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

func GiveMemoryTreeEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveMemoryTreeEnd: " + key)

	scrPath := filepath.Join(memoryTreeWorkingPath, key)
	workPath := filepath.Join(memoryTreeWorkingPath, key+".txt")
	unstagePath := filepath.Join(memoryTreeUstagePath, key+".txt")
	// unzip data
	err := file.DecompressionFile(scrPath, workPath, redis.RedisGetInt(key+"-MemoryTreeTotal"))
	if err != nil {
		return task.FAIL, err
	}
	// move to unstage
	err = file.MoveFile(workPath, unstagePath)
	if err != nil {
		return task.FAIL, err
	}
	// parse memory tree - todo
	content, err := os.ReadFile(unstagePath)
	if err != nil {
		return task.FAIL, err
	}
	err = handleRelation(content, key)
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

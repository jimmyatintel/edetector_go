package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var ExplorerMu *sync.Mutex
var UserExplorerChannel = make(map[string](chan string))

func init() {
	ExplorerMu = &sync.Mutex{}
}

func HandleExpolorer(p packet.Packet) {
	key := p.GetRkey()
	drives := strings.Split(p.GetMessage(), "|")
	redis.RedisSet(key+"-ExplorerProgress", 0)
	go updateDriveProgress(key)
	redis.RedisSet(key+"-DriveTotal", len(drives)-1)
	tmp_chan := make(chan string)
	ExplorerMu.Lock()
	UserExplorerChannel[key] = tmp_chan
	ExplorerMu.Unlock()
	for ind, d := range drives {
		parts := strings.Split(d, "-")
		if len(parts) == 2 {
			drive := parts[0]
			driveInfo := strings.Split(parts[1], ",")[0]
			msg := drive + "|" + driveInfo
			redis.RedisSet(key+"-DriveCount", ind)
			var user_packet = packet.TaskPacket{
				Key:     key,
				Message: msg,
			}
			err := StartGetExplorer(&user_packet)
			if err != nil {
				logger.Error("Start get explorer failed:", zap.Any("error", err.Error()))
			}
			ExplorerMu.Lock()
			block_chan := UserExplorerChannel[key]
			ExplorerMu.Unlock()
			block_chan <- msg
			logger.Info("Next round")
		}
	}
	logger.Info("Finish all drives: ", zap.Any("message", key))
	taskservice.Finish_task(key, "StartGetDrive")
}

func StartGetExplorer(p packet.UserPacket) error {
	err := clientsearchsend.SendUserTCPtoClient(p, task.EXPLORER_INFO, p.GetMessage())
	if err != nil {
		return err
	}
	return nil
}

package work

import (
	channelmap "edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"strings"

	"go.uber.org/zap"
)

func init() {
}

func HandleExpolorer(p packet.Packet) {
	key := p.GetRkey()
	drives := strings.Split(p.GetMessage(), "|")
	redis.RedisSet(key+"-ExplorerProgress", 0)
	go updateDriveProgress(key)
	redis.RedisSet(key+"-DriveTotal", len(drives)-1)
	tmp_chan := make(chan string)
	channelmap.AssignDiskChannel(key, &tmp_chan)
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
			block_chan, err := channelmap.GetDiskChannel(key)
			if err != nil {
				logger.Error("Get disk channel failed:", zap.Any("error", err.Error()))
				continue
			}
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

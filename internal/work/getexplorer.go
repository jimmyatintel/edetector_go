package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

var user_explorer = make(map[string]chan string)
var driveTotalMap = make(map[string]int)
var driveCountMap = make(map[string]int)

func HandleExpolorer(p packet.Packet) {
	key := p.GetRkey()
	drives := strings.Split(p.GetMessage(), "|")
	driveTotalMap[key] = len(drives) - 1
	go driveProgress(key)
	user_explorer[key] = make(chan string)
	for ind, d := range drives {
		parts := strings.Split(d, "-")
		if len(parts) == 2 {
			drive := parts[0]
			driveInfo := strings.Split(parts[1], ",")[0]
			msg := drive + "|" + driveInfo + "|Explorer|ScheduleName|0|2048"
			driveCountMap[key] = ind
			var user_packet = packet.TaskPacket{
				Key:     key,
				Message: msg,
			}
			err := StartGetExplorer(&user_packet)
			if err != nil {
				logger.Error("Start get explorer failed:", zap.Any("error", err.Error()))
			}
			m := strconv.Itoa(driveTotalMap[p.GetRkey()]) + "/" + strconv.Itoa(driveCountMap[p.GetRkey()]) + " " + msg
			logger.Info("Start handle & blocking ", zap.Any("message", m))
			user_explorer[key] <- msg
			logger.Info("Next round")
		}
	}
	taskservice.Finish_task(key, "StartGetDrive")
}

func StartGetExplorer(p packet.UserPacket) error {
	err := clientsearchsend.SendUserTCPtoClient(p, task.EXPLORER_INFO, p.GetMessage())
	if err != nil {
		return err
	}
	return nil
}

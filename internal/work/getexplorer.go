package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

var user_explorer = make(map[string]chan string)

func HandleExpolorer(p packet.Packet) {
	drives := strings.Split(p.GetMessage(), "|")
	user_explorer[p.GetRkey()] = make(chan string, 1)
	for _, d := range drives {
		parts := strings.Split(d, "-")
		if len(parts) == 2 {
			drive := parts[0]
			driveInfo := strings.Split(parts[1], ",")[0]
			msg := drive + "|" + driveInfo + "|Explorer|ScheduleName|0|2048"
			user_explorer[p.GetRkey()] <- msg
			fmt.Println("start handle ", msg)
			var user_packet = packet.TaskPacket{
				Key:     p.GetRkey(),
				Message: msg,
			}
			err := StartGetExplorer(&user_packet)
			if err != nil {
				logger.Error("Start get explorer failed:", zap.Any("error", err.Error()))
			}
		}
	}
	taskservice.Finish_task(p.GetRkey(), "StartGetDrive")
}

func StartGetExplorer(p packet.UserPacket) error {
	logger.Info("ExplorerInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.EXPLORER_INFO, p.GetMessage())
	if err != nil {
		return err
	}
	return nil
}

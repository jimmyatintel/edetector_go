package workfromapi

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"

	"go.uber.org/zap"
)

func Terminate(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("Terminate: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	// handlingTasks := taskservice.ListHandling_task(p.GetRkey())
	// for _, t := range handlingTasks {

	// }
	// err := clientsearchsend.SendUserTCPtoClient(p, task.GET_SCAN, p.GetMessage())
	// if err != nil {
	// 	return task.FAIL, err
	// }
	return task.SUCCESS, nil
}

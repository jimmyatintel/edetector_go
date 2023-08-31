package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"

	"go.uber.org/zap"
)

func Terminate(p packet.UserPacket) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("Terminate: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	err := clientsearchsend.SendUserTCPtoClient(p, task.TERMINATE_ALL, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	handlingTasks, err := query.Load_stored_task("nil", key, 2, "nil")
	if err != nil {
		return task.FAIL, err
	}
	for _, t := range handlingTasks {
		query.Terminated_task(t[0], key, t[3])
	}
	query.Finish_task(key, "Terminate")
	return task.SUCCESS, nil
}

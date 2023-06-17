package work

import (
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func CheckConnect(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	np := packet.CheckIsWork(p)
	logger.Info("CheckConnect: ", zap.Any("key", np.GetRKey()))
	// query.Update_time()
	return task.SUCCESS, nil
}

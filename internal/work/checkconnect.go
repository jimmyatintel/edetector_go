package work

import (
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func CheckConnect(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("CheckConnect: ", zap.Any("key", p.GetRkey()))
	// query.Update_time()
	return task.SUCCESS, nil
}

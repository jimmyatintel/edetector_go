package work

import (
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"
)

func CheckConnect(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("CheckConnect: " + p.GetRkey())
	// query.Update_time()
	return task.SUCCESS, nil
}

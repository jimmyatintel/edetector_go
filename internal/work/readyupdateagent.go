package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"net"

	"go.uber.org/zap"
)

var DataRight chan int

func ReadyUpdateAgent(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("ReadyUpdateAgent: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	path := filepath.Join("agentFile", "test.exe")
	fileInfo, err := os.Stat(path)
	if err != nil {
		return task.FAIL, err
	}
	fileLen := int(fileInfo.Size())
	err = clientsearchsend.SendTCPtoClient(p, task.GIVE_UPDATE_INFO, strconv.Itoa(fileLen), conn)
	if err != nil {
		return task.FAIL, err
	}
	go GiveUpdate(p, conn, fileLen, path)
	return task.SUCCESS, nil
}

func GiveUpdate(p packet.Packet, conn net.Conn, fileLen int, path string) {
	logger.Info("GiveUpdate: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Read file error", zap.Any("message", err.Error()))
	}
	start := 0
	for start < fileLen {
		<-DataRight
		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := content[start:end]
		err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_UPDATE, string(data), conn)
		if err != nil {
			logger.Error("Send GiveUpdate error", zap.Any("message", err.Error()))
		}
		start += 65436
	}
	err = clientsearchsend.SendDataTCPtoClient(p, task.GIVE_UPDATE_END, "", conn)
	if err != nil {
		logger.Error("Send GiveUpdateEnd error", zap.Any("message", err.Error()))
	}
}

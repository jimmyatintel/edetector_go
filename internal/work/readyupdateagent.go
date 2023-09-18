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
)

var DataRight chan net.Conn

func init() {
	DataRight = make(chan net.Conn)
}

func ReadyUpdateAgent(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("ReadyUpdateAgent: " + p.GetRkey() + "|" + p.GetMessage())
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
	go GiveUpdate(p, fileLen, path)
	return task.SUCCESS, nil
}

func GiveUpdate(p packet.Packet, fileLen int, path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Read file error: " + err.Error())
	}
	start := 0
	for {
		conn := <-DataRight
		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := content[start:end]
		logger.Info("GiveUpdate: " + p.GetRkey())
		err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_UPDATE, data, conn)
		if err != nil {
			logger.Error("Send GiveUpdate error: " + err.Error())
		}
		start += 65436
		if start >= fileLen {
			logger.Info("GiveUpdateEnd: " + p.GetRkey())
			err = clientsearchsend.SendTCPtoClient(p, task.GIVE_UPDATE_END, "", conn)
			if err != nil {
				logger.Error("Send GiveUpdateEnd error: " + err.Error())
			}
			break
		}
	}
}

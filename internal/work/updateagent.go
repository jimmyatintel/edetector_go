package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"net"
)

func ReadyUpdateAgent(p packet.Packet, conn net.Conn, dataRight chan net.Conn) (task.TaskResult, error) {
	logger.Info("ReadyUpdateAgent: " + p.GetRkey() + "::" + p.GetMessage())
	path := filepath.Join("agentFile", "test.exe")
	fileInfo, err := os.Stat(path)
	if err != nil {
		return task.FAIL, err
	}
	fileLen := int(fileInfo.Size())
	logger.Info("ServerSend GiveUpdateInfo: " + p.GetRkey() + "::" + strconv.Itoa(fileLen))
	err = clientsearchsend.SendTCPtoClient(p, task.GIVE_UPDATE_INFO, strconv.Itoa(fileLen), conn)
	if err != nil {
		return task.FAIL, err
	}
	go GiveUpdate(p, fileLen, path, dataRight)
	return task.SUCCESS, nil
}

func GiveUpdate(p packet.Packet, fileLen int, path string, dataRight chan net.Conn) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Read file error: " + err.Error())
	}
	start := 0
	for {
		conn := <-dataRight
		if start >= fileLen {
			logger.Info("ServerSend GiveUpdateEnd: " + p.GetRkey())
			err = clientsearchsend.SendDataTCPtoClient(p, task.GIVE_UPDATE_END, []byte{}, conn)
			if err != nil {
				logger.Error("Send GiveUpdateEnd error: " + err.Error())
			}
			<-dataRight
			query.Finish_task(p.GetRkey(), "StartUpdate")
			break
		}
		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := content[start:end]
		logger.Info("ServerSend GiveUpdate: " + p.GetRkey())
		err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_UPDATE, data, conn)
		if err != nil {
			logger.Error("Send GiveUpdate error: " + err.Error())
		}
		start += 65436
	}
}

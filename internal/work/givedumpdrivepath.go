package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/connectionmap"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"math"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func ReadyDumpDrive(p packet.Packet, conn net.Conn, dataRight chan net.Conn) (task.TaskResult, error) {
	key, path := p.GetRkey(), p.GetMessage()
	logger.Info(key + "::ReadyDumpDrive")

	// save dump drive path in ConnMsg
	redisKey := key + string(task.START_DUMP_DRIVE) + path
	taskId, err := redis.RedisGetString(redisKey)
	if err != nil {
		logger.Error("Error getting taskId from redis: " + err.Error())
		return task.FAIL, err
	}

	connectionmap.StoreConnInfo(conn, connectionmap.ConnInfo{
		TaskId:  taskId,
		Msg:     path,
		DataLen: 0,
	})

	go GiveDumpDrivePathInfo(p, conn, dataRight)

	return task.SUCCESS, nil
}

func GiveDumpDrivePathInfo(p packet.Packet, conn net.Conn, dataRight chan net.Conn) {
	logger.Info(p.GetRkey() + "::GiveDumpDrivePathInfo")

	// read dump drive path from ConnMsg
	connInfo, ok := connectionmap.GetConnInfo(conn)
	if !ok {
		logger.Error("Error getting msg from ConnMsgMap")
		return
	}

	// read dump drive file size
	readPath := filepath.Join(dumpWorkingPath, connInfo.TaskId+".txt")
	content, err := os.ReadFile(readPath)
	if err != nil {
		logger.Error("GetFileSize: " + err.Error())
		return
	}
	fileSizeStr := strconv.Itoa(len(content))

	// send file size to client
	err = clientsearchsend.SendTCPtoClient(p, task.GIVE_DUMP_DRIVE_INFO, fileSizeStr, conn)
	if err != nil {
		logger.Error("SendTCPtoClient: " + err.Error())
		return
	}

	go GiveDumpDrivePath(p, len(content), string(content), dataRight)

	// remove txt in working path
	if err := os.Remove(filepath.Join(dumpWorkingPath, connInfo.TaskId+".txt")); err != nil {
		logger.Error("Error removing txt in working path: " + err.Error())
	}
}

func GiveDumpDrivePath(p packet.Packet, fileLen int, content string, dataRight chan net.Conn) {
	start := 0
	for {
		conn := <-dataRight
		if start >= fileLen {
			logger.Info(p.GetRkey() + "::GiveDumpDrivePathEnd")
			err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_DUMP_DRIVE_END, []byte{}, conn)
			if err != nil {
				logger.Error("SendDataTCPtoClient: " + err.Error())
				return
			}
			<-dataRight
			break
		}

		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := []byte(content[start:end])
		logger.Info(p.GetRkey() + "::GiveDumpDrivePath: " + strconv.Itoa(start) + "-" + strconv.Itoa(end))

		err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_DUMP_DRIVE_PATH, data, conn)
		if err != nil {
			logger.Error("SendDataTCPtoClient: " + err.Error())
			return
		}

		start += 65436
	}
}

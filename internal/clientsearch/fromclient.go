package clientsearch

import (
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/task"
	"fmt"
	"strings"

	channelmap "edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	work "edetector_go/internal/work"
	logger "edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"net"
)

var Clientlist []string

func handleTCPRequest(conn net.Conn, task_chan chan packet.Packet, port string) {
	defer conn.Close()
	buf := make([]byte, 2048)
	key := "unknown"
	agentTaskType := "unknown"
	lastTask := "unknown"
	var dataRightChan chan net.Conn
	closeConn := make(chan bool)
	for {
		var NewPacket packet.Packet
		var decrypt_buf []byte
		reqLen, err := conn.Read(buf)
		// debug
		// var temp_buf []byte
		// logger.Debug("Read len: " + strconv.Itoa(reqLen))
		// temp_buf = bytes.Repeat([]byte{0}, reqLen)
		// C_AES.Decryptbuffer(buf, reqLen, temp_buf)
		// logger.Debug("Read tmp buffer: " + string(temp_buf))
		// debug end
		if err != nil {
			connectionClosedByAgent(key, agentTaskType, lastTask, err)
			close(closeConn)
			return
		}
		Data_acache := make([]byte, 0)
		Data_acache = append(Data_acache, buf[:reqLen]...)
		if reqLen == 1024 {
			NewPacket = new(packet.WorkPacket)
		} else if reqLen > 1024 {
			for len(Data_acache) < 65535 {
				reqLen, err := conn.Read(buf)
				if err != nil {
					connectionClosedByAgent(key, agentTaskType, lastTask, err)
					close(closeConn)
					return
				}
				// debug
				// var temp_buf []byte
				// logger.Debug("Read len: " + strconv.Itoa(reqLen))
				// temp_buf = bytes.Repeat([]byte{0}, reqLen)
				// C_AES.Decryptbuffer(buf[:reqLen], reqLen, temp_buf)
				// logger.Debug("Read tmp buffer: " + string(temp_buf))
				// debug end
				Data_acache = append(Data_acache, buf[:reqLen]...)
			}
			NewPacket = new(packet.DataPacket)
		} else {
			logger.Error("Invalid packet (too short): " + fmt.Sprintf("%x", decrypt_buf))
			continue
		}
		decrypt_buf = bytes.Repeat([]byte{0}, len(Data_acache))
		C_AES.Decryptbuffer(Data_acache, len(Data_acache), decrypt_buf)
		err = NewPacket.NewPacket(decrypt_buf, Data_acache)
		if err != nil {
			logger.Error("Error reading: " + err.Error() + ", Content: " + fmt.Sprintf("%x", decrypt_buf))
			continue
		}
		if key == "unknown" {
			key = NewPacket.GetRkey()
		}
		if agentTaskType == "unknown" {
			taskType, ok := task.TaskTypeMap[NewPacket.GetTaskType()]
			if ok {
				agentTaskType = taskType
			}
		}
		if NewPacket.GetTaskType() == "Undefine" {
			nullIndex := bytes.IndexByte(decrypt_buf[76:100], 0)
			logger.Error("Undefine TaskType: " + string(decrypt_buf[76:76+nullIndex]))
			continue
		}
		if NewPacket.GetTaskType() == task.GIVE_INFO {
			go func() {
				for {
					select {
					case message := <-task_chan:
						data := message.Fluent()
						logger.Info("Get task msg: " + string(data))
						err := clientsearchsend.SendTaskTCPtoClient(data, conn)
						if err != nil {
							logger.Error("Error Sending: " + err.Error())
						}
					case <-closeConn:
						return
					}
				}
			}()
			// wait for key to join the packet
			Clientlist = append(Clientlist, key)
			channelmap.AssignTaskChannel(key, &task_chan)
			logger.Info("Set key-channel mapping: " + key)
		} else {
			redis.Online(key)
		}
		if NewPacket.GetTaskType() == task.READY_UPDATE_AGENT {
			dataRightChan = make(chan net.Conn)
			work.ReadyUpdateAgent(NewPacket, conn, dataRightChan)
		} else if agentTaskType == "StartUpdate" && NewPacket.GetTaskType() == task.DATA_RIGHT {
			logger.Info("DataRight: " + key)
			dataRightChan <- conn
		} else {
			taskFunc, ok := work.WorkMap[NewPacket.GetTaskType()]
			if !ok {
				logger.Error("Function notfound: " + string(NewPacket.GetTaskType()))
				continue
			}
			_, err = taskFunc(NewPacket, conn)
			if err != nil {
				logger.Error("Task " + string(NewPacket.GetTaskType()) + " failed: " + err.Error())
				if agentTaskType == "StartScan" || agentTaskType == "StartGetDrive" || agentTaskType == "StartCollect" || agentTaskType == "StartGetImage" {
					query.Failed_task(NewPacket.GetRkey(), agentTaskType)
				}
				continue
			}
		}
		lastTask = string(NewPacket.GetTaskType())
	}
}

func handleUDPRequest(addr net.Addr, buf []byte) {
	logger.Info("UDP")
}

func connectionClosedByAgent(key string, agentTaskType string, lastTask string, err error) {
	if agentTaskType != "CollectProgress" {
		logger.Warn("Connection close: " + string(key) + "|" + agentTaskType + ", Error: " + err.Error())
	}
	if agentTaskType == "StartScan" && lastTask == "ReadyScan" {
		query.Update_task_status(key, agentTaskType, 2, 0)
	} else if agentTaskType == "StartScan" || agentTaskType == "StartGetDrive" || agentTaskType == "StartCollect" || agentTaskType == "StartGetImage" {
		if !strings.Contains(lastTask, "End") {
			query.Failed_task(key, agentTaskType)
		}
	} else if agentTaskType == "Main" {
		redis.Offline(key, false)
	}
}

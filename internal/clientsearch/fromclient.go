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
	logger.Info("Worker port accepted, IP: " + conn.RemoteAddr().String())
	defer conn.Close()
	buf := make([]byte, 1024)
	key := "unknown"
	agentTaskType := "unknown"
	lastTask := "unknown"
	var dataRightChan chan net.Conn
	closeConn := make(chan bool)
	for {
		var NewPacket packet.Packet
		var decrypt_buf []byte
		reqLen, err := conn.Read(buf)
		if err != nil {
			connectionClosedByAgent(key, agentTaskType, lastTask, err)
			close(closeConn)
			return
		}
		if reqLen < 1024 {
			logger.Error("Invalid packet (too short): " + fmt.Sprintf("%x", buf[:reqLen]))
			continue
		}
		Data_acache := make([]byte, 0)
		Data_acache = append(Data_acache, buf[:reqLen]...)
		decrypt_buf = bytes.Repeat([]byte{0}, len(Data_acache))
		C_AES.Decryptbuffer(Data_acache, len(Data_acache), decrypt_buf)
		NewPacket = new(packet.WorkPacket)
		err = NewPacket.NewPacket(decrypt_buf, Data_acache)
		if err != nil {
			logger.Error("Error reading: " + err.Error())
			logger.Info("Content: " + fmt.Sprintf("%x", decrypt_buf))
			continue
		}
		t := NewPacket.GetTaskType()
		if t != "GiveInfo" && t != "GiveDetectInfoFirst" && t != "GiveDetectInfo" && t != "CheckConnect" && t != "ReadyUpdateAgent" && t != "DataRight" && t != "GiveDriveInfo" {
			for len(Data_acache) < 65535 {
				reqLen, err := conn.Read(buf)
				if err != nil {
					connectionClosedByAgent(key, agentTaskType, lastTask, err)
					close(closeConn)
					return
				}
				Data_acache = append(Data_acache, buf[:reqLen]...)
			}
			NewPacket = new(packet.DataPacket)
			decrypt_buf = bytes.Repeat([]byte{0}, len(Data_acache))
			C_AES.Decryptbuffer(Data_acache, len(Data_acache), decrypt_buf)
			// rand := fmt.Sprint(rand.Intn(256))
			// tmpPath := "test/encrypted_long_" + rand
			// file.WriteFile(tmpPath, Data_acache)
			// tmpPath = "test/decrypted_long_" + rand
			// file.WriteFile(tmpPath, decrypt_buf)
			err = NewPacket.NewPacket(decrypt_buf, Data_acache)
			if err != nil {
				logger.Error("Error reading: " + err.Error())
				logger.Info("Content: " + fmt.Sprintf("%x", decrypt_buf))
				continue
			}
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

// To-Do (TBD)
func connectionClosedByAgent(key string, agentTaskType string, lastTask string, err error) {
	if agentTaskType == "StartScan" && lastTask == "ReadyScan" {
		logger.Warn("Connection close: " + string(key) + "|" + agentTaskType + ", Error: " + err.Error())
		query.Update_task_status(key, agentTaskType, 2, 0)
	} else if agentTaskType == "StartScan" || agentTaskType == "StartGetDrive" || agentTaskType == "StartCollect" || agentTaskType == "StartGetImage" {
		if !strings.Contains(lastTask, "End") {
			logger.Warn("Connection close: " + string(key) + "|" + agentTaskType + ", Error: " + err.Error())
			query.Failed_task(key, agentTaskType)
		}
	} else if agentTaskType == "Main" {
		logger.Warn("Connection close: " + string(key) + "|" + agentTaskType + ", Error: " + err.Error())
		removeTasks, err := query.Load_stored_task("nil", key, 2, "StartRemove")
		if err != nil {
			logger.Error("Get StartRemove tasks failed: " + err.Error())
		}
		if len(removeTasks) == 0 {
			redis.Offline(key, false)
		} else {
			query.DeleteAgent(key)
			err = redis.RedisDelete(key)
			if err != nil {
				logger.Error("Error deleting key from redis: " + err.Error())
			}
			logger.Info("Finish remove agent: " + key)
		}
	}
}

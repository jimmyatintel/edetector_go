package clientsearch

import (
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/task"

	channelmap "edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	work "edetector_go/internal/work"
	logger "edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"net"

	"go.uber.org/zap"
)

var Clientlist []string

func handleTCPRequest(conn net.Conn, task_chan chan packet.Packet, port string) {
	defer conn.Close()
	buf := make([]byte, 2048)
	key := "unknown"
	agentTaskType := "unknown"
	retryScanFlag := false
	if task_chan != nil {
		go func() {
			for {
				select {
				case message := <-task_chan:
					data := message.Fluent()
					logger.Info("get task msg: ", zap.Any("message", string(data)))
					err := clientsearchsend.SendTaskTCPtoClient(data, conn)
					if err != nil {
						logger.Error("send failed:", zap.Any("error", err.Error()))
					}
				}
			}
		}()
	}
	for {
		reqLen, err := conn.Read(buf)
		if err != nil {
			if agentTaskType != "CollectProgress" {
				logger.Warn("Connection close: ", zap.Any("message", string(key)+" ,Type: "+agentTaskType+" ,Error: "+err.Error()))
			}
			if agentTaskType == "StartScan" && retryScanFlag {
				query.Update_task_status(key, agentTaskType, 2, 0)
			} else if agentTaskType == "StartScan" || agentTaskType == "StartGetDrive" || agentTaskType == "StartCollect" {
				query.Failed_task(key, agentTaskType)
			} else if agentTaskType == "Main" {
				redis.Offline(key)
			}
			return
		}
		decrypt_buf := bytes.Repeat([]byte{0}, reqLen)
		C_AES.Decryptbuffer(buf, reqLen, decrypt_buf)
		if reqLen == 1024 {
			var NewPacket = new(packet.WorkPacket)
			err := NewPacket.NewPacket(decrypt_buf, buf)
			if err != nil {
				logger.Error("Error reading:", zap.Any("error", err.Error()+" "+string(decrypt_buf)))
				return
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
				logger.Error("Undefine Task Type: ", zap.String("error", string(decrypt_buf[76:76+nullIndex])))
				logger.Error("pkt content: ", zap.String("error", string(NewPacket.GetMessage())))
				return
			}
			if NewPacket.GetTaskType() == task.READY_SCAN {
				retryScanFlag = true
			} else {
				retryScanFlag = false
			}
			if NewPacket.GetTaskType() == task.GIVE_INFO {
				// wait for key to join the packet
				Clientlist = append(Clientlist, key)
				channelmap.AssignTaskChannel(key, &task_chan)
				logger.Info("set worker key-channel mapping: ", zap.Any("message", key))
			} else {
				redis.Online(key)
			}
			taskFunc, ok := work.WorkMap[NewPacket.GetTaskType()]
			if !ok {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()))
				return
			}
			_, err = taskFunc(NewPacket, conn)
			if err != nil {
				logger.Error(string(NewPacket.GetTaskType())+" task failed: ", zap.Any("error", err.Error()))
				if agentTaskType == "StartScan" || agentTaskType == "StartGetDrive" || agentTaskType == "StartCollect" {
					query.Failed_task(NewPacket.GetRkey(), agentTaskType)
				} else if agentTaskType == "Main" {
					redis.Offline(key)
				}
				return
			}
		} else if reqLen > 1024 {
			Data_acache := make([]byte, 0)
			Data_acache = append(Data_acache, buf[:reqLen]...)
			for len(Data_acache) < 65535 {
				reqLen, err := conn.Read(buf)
				if err != nil {
					logger.Warn("Connection close: ", zap.Any("message", string(key)+" ,Type: "+agentTaskType+" ,Error: "+err.Error()))
					if agentTaskType == "StartScan" || agentTaskType == "StartGetDrive" || agentTaskType == "StartCollect" {
						query.Failed_task(key, agentTaskType)
					}
					return
				}
				Data_acache = append(Data_acache, buf[:reqLen]...)
			}
			decrypt_buf := bytes.Repeat([]byte{0}, len(Data_acache))
			C_AES.Decryptbuffer(Data_acache, len(Data_acache), decrypt_buf)
			var NewPacket = new(packet.DataPacket)
			err := NewPacket.NewPacket(decrypt_buf, Data_acache)
			if err != nil {
				logger.Error("Error reading:", zap.Any("error", err.Error()+" "+string(decrypt_buf)))
				return
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
				logger.Error("Undefine Task Type: ", zap.String("error", string(decrypt_buf[76:76+nullIndex])))
				logger.Error("pkt content: ", zap.String("error", string(NewPacket.GetMessage())))
				return
			}
			retryScanFlag = false
			redis.Online(key)
			taskFunc, ok := work.WorkMap[NewPacket.GetTaskType()]
			if !ok {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()))
				return
			}
			_, err = taskFunc(NewPacket, conn)
			if err != nil {
				logger.Error(string(NewPacket.GetTaskType())+" task failed:", zap.Any("error", err.Error()))
				if agentTaskType == "StartScan" || agentTaskType == "StartGetDrive" || agentTaskType == "StartCollect" {
					query.Failed_task(key, agentTaskType)
				}
				return
			}
		} else {
			logger.Error("Invalid packet(short):", zap.Any("message", decrypt_buf))
			continue
		}
	}
}

func handleUDPRequest(addr net.Addr, buf []byte) {
	logger.Info("udp")
}

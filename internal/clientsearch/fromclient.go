package clientsearch

import (
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/task"
	"edetector_go/internal/taskservice"

	channelmap "edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	work "edetector_go/internal/work"
	logger "edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
	"net"

	"go.uber.org/zap"
)

var Clientlist []string

func handleTCPRequest(conn net.Conn, task_chan chan packet.Packet, port string) {
	defer conn.Close()
	buf := make([]byte, 2048)
	key := "null"
	if task_chan != nil {
		go func() {
			for {
				select {
				case message := <-task_chan:
					data := message.Fluent()
					logger.Info("get task msg: ", zap.Any("message", string(data)))
					err := clientsearchsend.SendTCPtoClient(data, conn)
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
			if err.Error() == "EOF" {
				if key != "null" {
					err = redis.Offline(key)
					if err != nil {
						logger.Error("Update offline error:", zap.Any("error", err.Error()))
					}
				}
				logger.Info(string(key) + " " + string(port) + " Connection close")
				return
			} else {
				logger.Error("Error reading:", zap.Any("error", string(key)+err.Error()))
				return
			}
		}
		decrypt_buf := bytes.Repeat([]byte{0}, reqLen)
		C_AES.Decryptbuffer(buf, reqLen, decrypt_buf)
		// fmt.Println("decrypt buf: ", string(decrypt_buf))
		if reqLen == 1024 {
			rabbitmq.Declare("clientsearch")
			var NewPacket = new(packet.WorkPacket)
			err := NewPacket.NewPacket(decrypt_buf, buf)
			if err != nil {
				logger.Error("Error reading:", zap.Any("error", err.Error()+" "+string(decrypt_buf)))
				return
			}
			if NewPacket.GetTaskType() == "Undefine" {
				nullIndex := bytes.IndexByte(decrypt_buf[76:100], 0)
				logger.Error("Undefine Task Type: ", zap.String("error", string(decrypt_buf[76:76+nullIndex])))
				logger.Error("pkt content: ", zap.String("error", string(NewPacket.GetMessage())))
				return
			}
			// fmt.Println("Get task: ", NewPacket.GetTaskType())
			if NewPacket.GetTaskType() == task.GIVE_INFO {
				// wait for key to join the packet
				key = NewPacket.GetRkey()
				Clientlist = append(Clientlist, key)
				channelmap.AssignTaskChannel(key, &task_chan)
				logger.Info("set worker key-channel mapping: ", zap.Any("message", NewPacket.GetRkey()))
			} else if key != "null" {
				err = redis.Online(key)
				if err != nil {
					logger.Error("Update online failed:", zap.Any("error", err.Error()))
				}
				if NewPacket.GetTaskType() == task.GIVE_DETECT_INFO_FIRST {
					taskservice.RequestToUser(key)
				}
			}
			taskFunc, ok := work.WorkMap[NewPacket.GetTaskType()]
			if !ok {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()))
				return
			}
			_, err = taskFunc(NewPacket, conn)
			if err != nil {
				logger.Error(string(NewPacket.GetTaskType())+" task failed: ", zap.Any("error", err.Error()))
				taskType, ok := task.TaskTypeMap[NewPacket.GetTaskType()]
				if ok {
					taskservice.Failed_task(NewPacket.GetRkey(), taskType)
				}
				return
			}
		} else if reqLen > 1024 {
			Data_acache := make([]byte, 0)
			Data_acache = append(Data_acache, buf[:reqLen]...)
			for len(Data_acache) < 65535 {
				reqLen, err := conn.Read(buf)
				if err != nil {
					if err.Error() == "EOF" {
						if key != "null" {
							redis.Offline(key)
						}
						logger.Info(string(key) + " " + string(port) + " Connection close")
						return
					} else {
						logger.Error("Error reading:", zap.Any("error", err.Error()))
						return
					}
				}
				Data_acache = append(Data_acache, buf[:reqLen]...)
			}
			decrypt_buf := bytes.Repeat([]byte{0}, len(Data_acache))
			C_AES.Decryptbuffer(Data_acache, len(Data_acache), decrypt_buf)
			// logger.Info("Receive Large TCP from client", zap.Any("data", string(decrypt_buf)), zap.Any("len", len(Data_acache)))
			var NewPacket = new(packet.DataPacket)
			err := NewPacket.NewPacket(decrypt_buf, Data_acache)
			if err != nil {
				logger.Error("Error reading:", zap.Any("error", err.Error()+" "+string(decrypt_buf)))
				return
			}
			if NewPacket.GetTaskType() == "Undefine" {
				nullIndex := bytes.IndexByte(decrypt_buf[76:100], 0)
				logger.Error("Undefine Task Type: ", zap.String("error", string(decrypt_buf[76:76+nullIndex])))
				logger.Error("pkt content: ", zap.String("error", string(NewPacket.GetMessage())))
				return
			}
			// fmt.Println("task type: ", NewPacket.GetTaskType(), port)
			if NewPacket.GetTaskType() != task.GIVE_INFO && NewPacket.GetTaskType() != task.GIVE_DETECT_INFO_FIRST {
				err = redis.Online(key)
				if err != nil {
					logger.Error("Upate online failed:", zap.Any("error", err.Error()))
				}
			}
			taskFunc, ok := work.WorkMap[NewPacket.GetTaskType()]
			if !ok {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()))
				return
			}
			_, err = taskFunc(NewPacket, conn)
			if err != nil {
				logger.Error(string(NewPacket.GetTaskType())+" task failed:", zap.Any("error", err.Error()))
				taskType, ok := task.TaskTypeMap[NewPacket.GetTaskType()]
				if ok {
					taskservice.Failed_task(NewPacket.GetRkey(), taskType)
				}
				return
			}
		} else {
			logger.Error("Invalid packet(short):", zap.Any("message", decrypt_buf))
		}
	}
}

func handleUDPRequest(addr net.Addr, buf []byte) {
	logger.Info("udp")
}

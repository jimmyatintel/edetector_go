package clientsearch

import (
	// "context"
	// config "edetector_go/config"
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/task"

	// fflag "edetector_go/internal/fflag"
	packet "edetector_go/internal/packet"
	work "edetector_go/internal/work"
	work_from_api "edetector_go/internal/work_from_api"
	logger "edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"fmt"

	"go.uber.org/zap"

	// "io/ioutil"
	"net"
)

var Key *string

func handleTCPRequest(conn net.Conn, task_chan chan string) {
	defer conn.Close()
	buf := make([]byte, 2048)
	Key = new(string)
	*Key = "null"
	defer func() {
		Key = nil
	}()
	if task_chan != nil { //! 不確定具體用途
		go func() {
			for {
				select {
				case message := <-task_chan:
					fmt.Println("get task msg: " + message)
				}
			}
		}()
	}
	for {
		reqLen, err := conn.Read(buf)
		// buf, err := ioutil.ReadAll(conn)
		// reqLen := len(buf)
		if err != nil {
			if err.Error() == "EOF" {
				logger.Debug("Connection close")
				return
			} else {
				logger.Error("Error reading:", zap.Any("error", err.Error()))
				return
			}
		}
		decrypt_buf := bytes.Repeat([]byte{0}, reqLen)
		C_AES.Decryptbuffer(buf, reqLen, decrypt_buf)

		if reqLen <= 1024 && Key != nil {
			// fmt.Println(string(decrypt_buf))
			rabbitmq.Declare("clientsearch")
			var NewPacket = new(packet.WorkPacket)
			err := NewPacket.NewPacket(decrypt_buf, buf)
			if err != nil {
				logger.Error("Error reading:", zap.Any("error", err.Error()), zap.Any("len", reqLen))
				return
			}
			logger.Info("Receive TCP from client", zap.Any("function", NewPacket.GetTaskType()))
			_, err = work.WorkMap[NewPacket.GetTaskType()](NewPacket, Key, conn)
			if NewPacket.GetTaskType() == task.GIVE_INFO {
				fmt.Println(*Key)
			}
			if err != nil {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()), zap.Any("error", err.Error()))
				return
			}
		} else if reqLen > 0 && Key != nil && *Key == "null" {
			Data_acache := make([]byte, 0)
			Data_acache = append(Data_acache, buf[:reqLen]...)
			for len(Data_acache) < 65535 {
				reqLen, err := conn.Read(buf)
				if err != nil {
					if err.Error() == "EOF" {
						logger.Debug("Connection close")
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
				logger.Info("Receive TCP from client", zap.Any("function", NewPacket.GetTaskType()))
				logger.Error("Error reading:", zap.Any("error", err.Error()), zap.Any("len", reqLen))
				return
			}
			_, err = work.WorkMap[NewPacket.GetTaskType()](NewPacket, Key, conn)
			if err != nil {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()))
				return
			}
		}
		// wait for key to join the packet
		if *Key != "null" && task_chan != nil {
			Task_channel[*Key] = task_chan
			fmt.Println("set task " + *Key)
		}
	}
}

func handleTaskrequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	Key = new(string)
	*Key = "null"
	defer func() {
		Key = nil
	}()

	for {
		reqLen, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				logger.Debug("Connection close")
				return
			} else {
				logger.Error("Error reading:", zap.Any("error", err.Error()))
				return
			}
		}
		if reqLen <= 1024 {
			content := buf[:reqLen]
			NewPacket := new(packet.TaskPacket)
			err = NewPacket.NewPacket(content)
			if err != nil {
				logger.Error("Error reading task packet:", zap.Any("error", err.Error()), zap.Any("len", reqLen))
				return
			}
			logger.Info("Receive task from user", zap.Any("function", NewPacket.GetUserTaskType()))
			_, err = work_from_api.WorkapiMap[NewPacket.GetUserTaskType()](NewPacket, Key, conn)
			if err != nil {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetUserTaskType()), zap.Any("error", err.Error()))
				return
			}
		}else {
			logger.Error("Task packet is longer than 1024")
			return
		}
	}
}
func handleUDPRequest(addr net.Addr, buf []byte) {
	fmt.Println(string(buf))
}

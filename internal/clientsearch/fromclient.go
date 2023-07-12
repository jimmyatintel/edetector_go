package clientsearch

import (
	// "context"
	// config "edetector_go/config"
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/task"

	// fflag "edetector_go/internal/fflag"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	taskchannel "edetector_go/internal/taskchannel"
	work "edetector_go/internal/work"
	logger "edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"fmt"

	"go.uber.org/zap"

	// "io/ioutil"
	"net"
	// "encoding/binary"
)


func handleTCPRequest(conn net.Conn, task_chan chan packet.Packet, port string) {
	defer conn.Close()
	buf := make([]byte, 2048)
	if task_chan != nil {
		go func() {
			for {
				select {
				case message := <-task_chan:
					data := message.Fluent()
					fmt.Println(len(data))
					fmt.Println("get task msg: " + string(data))
					err := clientsearchsend.SendTCPtoClient(data, conn)
					if err != nil {
						logger.Error("Send failed:", zap.Any("error", err.Error()))
					}
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
		if reqLen <= 1024 {
			// fmt.Println(string(decrypt_buf))
			rabbitmq.Declare("clientsearch")
			var NewPacket = new(packet.WorkPacket)
			err := NewPacket.NewPacket(decrypt_buf, buf)
			if err != nil {
				logger.Error("Error reading:", zap.Any("error", err.Error()), zap.Any("len", reqLen))
				return
			}
			if NewPacket.GetTaskType() == "Undefine" {
				nullIndex := bytes.IndexByte(decrypt_buf[76:100], 0)
				logger.Error("Undefine Task Type: ", zap.String("error", string(decrypt_buf[76:76+nullIndex])))
				logger.Error("pkt content: ", zap.String("error", string(NewPacket.GetMessage())))
				return
			}
			// fmt.Println("task type: ", NewPacket.GetTaskType(), port)
			taskFunc, ok := work.WorkMap[NewPacket.GetTaskType()]
			if !ok {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()))
				return
			}
			_, err = taskFunc(NewPacket, conn)
			if err != nil {
				logger.Error("Task Failed:", zap.Any("error", err.Error()))
				return
			}
			if NewPacket.GetTaskType() == task.GIVE_INFO {
				// wait for key to join the packet
				taskchannel.Task_worker_channel[NewPacket.GetRkey()] = task_chan
				fmt.Println("set worker key-channel mapping: " + NewPacket.GetRkey())
			}
		} else if reqLen > 0 {
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
				logger.Error("Error reading:", zap.Any("error", err.Error()), zap.Any("len", reqLen))
				return
			}
			if NewPacket.GetTaskType() == "Undefine" {
				nullIndex := bytes.IndexByte(decrypt_buf[76:100], 0)
				logger.Error("Undefine Task Type: ", zap.String("error", string(decrypt_buf[76:76+nullIndex])))
				logger.Error("pkt content: ", zap.String("error", string(NewPacket.GetMessage())))
				return
			}
			// fmt.Println("task type: ", NewPacket.GetTaskType(), port)
			taskFunc, ok := work.WorkMap[NewPacket.GetTaskType()]
			if !ok {
				logger.Error("Function notfound:", zap.Any("name", NewPacket.GetTaskType()))
				return
			}
			_, err = taskFunc(NewPacket, conn)
			if err != nil {
				logger.Error("Task Failed:", zap.Any("error", err.Error()))
				return
			}
		}
	}
}

func handleUDPRequest(addr net.Addr, buf []byte) {
	fmt.Println(string(buf))

}

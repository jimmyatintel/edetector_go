package clientsearchsend

import (
	// "context"
	C_AES "edetector_go/internal/C_AES"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/internal/taskchannel"
	"fmt"

	"errors"
	"net"
)

func init() {

}

func SendTCPtoClient(data []byte, conn net.Conn) error {
	encrypt_buf := make([]byte, len(data))
	// logger.Info("Send TCP to client", zap.Any("len", len(data)))
	C_AES.Encryptbuffer(data, len(data), encrypt_buf)
	_, err := conn.Write(encrypt_buf)
	if err != nil {
		return err
	}
	return nil
}

func SendUserTCPtoClient(p packet.UserPacket, workType task.TaskType, msg string, port string) error{
	fmt.Println("port", port)
	if port == "worker" {
		var send_packet = packet.WorkPacket{
			MacAddress: p.GetMacAddress(),
			IpAddress:  p.GetipAddress(),
			Work:       workType,
			Message:    msg,
		}
		task_chan := taskchannel.Task_worker_channel[p.GetRkey()]
		task_chan <- &send_packet
		return nil
	}
	return errors.New("invalid port")
}

func SendDriveTCPtoClient(p packet.Packet, key string, workType task.TaskType, msg string) error{
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       workType,
		Message:    msg,
	}
	task_chan := taskchannel.Task_worker_channel[key]
	task_chan <- &send_packet
	return nil
}
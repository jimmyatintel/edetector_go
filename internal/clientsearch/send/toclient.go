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
	// var send_packet = packet.WorkPacket{
	// 	MacAddress: p.GetMacAddress(),
	// 	IpAddress:  p.GetipAddress(),
	// 	Work:       workType,
	// 	Message:    msg, // "0|0"
	// }
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
	} else if port == "detect" {
		var send_packet = packet.DataPacket{
			MacAddress: p.GetMacAddress(),
			IpAddress:  p.GetipAddress(),
			Work:       workType,
			Message:    msg,
		}
		fmt.Println("detect!!")
		task_chan := taskchannel.Task_detect_channel[p.GetRkey()]
		task_chan <- &send_packet
		return nil
	}
	return errors.New("invalid port")
}

// func SendDetectTCPtoClient(p packet.Packet, workType task.TaskType, msg string, port string) error{
// 	fmt.Println("port", port)

// 	if port == "worker" {
// 		var send_packet = packet.WorkPacket{
// 			MacAddress: p.GetMacAddress(),
// 			IpAddress:  p.GetipAddress(),
// 			Work:       workType,
// 			Message:    msg,
// 		}
// 		task_chan := taskchannel.Task_worker_channel[p.GetRkey()]
// 		task_chan <- send_packet.Fluent()
// 		return nil
// 	} else if port == "detect" {
// 		var send_packet = packet.DataPacket{
// 			MacAddress: p.GetMacAddress(),
// 			IpAddress:  p.GetipAddress(),
// 			Work:       workType,
// 			Message:    msg,
// 		}
// 		task_chan := taskchannel.Task_detect_channel[p.GetRkey()]
// 		task_chan <- send_packet.Fluent()
// 		return nil
// 	}
// 	return errors.New("invalid port")
// }
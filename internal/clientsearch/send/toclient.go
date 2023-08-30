package clientsearchsend

import (
	// "context"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/channelmap"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"

	"net"
)

func SendTaskTCPtoClient(data []byte, conn net.Conn) error {
	encrypt_buf := make([]byte, len(data))
	C_AES.Encryptbuffer(data, len(data), encrypt_buf)
	_, err := conn.Write(encrypt_buf)
	if err != nil {
		return err
	}
	return nil
}

func SendTCPtoClient(p packet.Packet, worktype task.TaskType, msg string, conn net.Conn) error {
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       worktype,
		Message:    msg,
	}
	data := send_packet.Fluent()
	encrypt_buf := make([]byte, len(data))
	C_AES.Encryptbuffer(data, len(data), encrypt_buf)
	_, err := conn.Write(encrypt_buf)
	if err != nil {
		return err
	}
	return nil
}

func SendUserTCPtoClient(p packet.UserPacket, workType task.TaskType, msg string) error {
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       workType,
		Message:    msg,
	}
	task_chan, err := channelmap.GetTaskChannel(p.GetRkey())
	if err != nil {
		return err
	}
	task_chan <- &send_packet
	return nil
}

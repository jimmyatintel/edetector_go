package clientsearchsend

import (
	// "context"s
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/channelmap"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"strconv"
	"strings"

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

func appendByteMsg(data []byte, msg []byte) []byte {
	length := 65436
	if len(msg) > length {
		logger.Error("Error msg length of DataPacket")
		return data
	} else if len(msg) == length {
		msg = msg[:length]
	} else {
		msg = append(msg, []byte(strings.Repeat(string(" "), length-len(msg)))...)
	}
	data = data[:100]
	data = append(data, msg...)
	return data
}

func SendDataTCPtoClient(p packet.Packet, worktype task.TaskType, msg []byte, conn net.Conn) error {
	var send_packet = packet.DataPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       worktype,
		Message:    "",
	}
	data := send_packet.Fluent()
	data = appendByteMsg(data, msg)
	encrypt_buf := make([]byte, len(data))
	logger.Info("Data len: " + strconv.Itoa(len(data)))
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

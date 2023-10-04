package mockagent

import (
	"edetector_go/internal/C_AES"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"
)

func SendTCPtoServer(worktype task.TaskType, msg string, conn net.Conn, ip string, mac string, key string) error {
	logger.Info(key + ":: SendTCPtoServer: " + string(worktype) + " " + msg)
	var send_packet = packet.WorkPacket{
		MacAddress: ip,
		IpAddress:  mac,
		Rkey:       key,
		Work:       worktype,
		Message:    msg,
	}
	data := send_packet.Fluent()
	encrypt_buf := make([]byte, len(data))
	C_AES.Encryptbuffer(data, len(data), encrypt_buf)
	_, err := conn.Write(encrypt_buf)
	if err != nil {
		logger.Error(key + ":: Error sending data to server: " + err.Error())
	}
	return nil
}

func SendDataTCPtoServer(worktype task.TaskType, msg []byte, conn net.Conn, ip string, mac string, key string) error {
	// logger.Debug(key+":: SendDataTCPtoServer: " + string(worktype))
	var send_packet = packet.DataPacket{
		MacAddress: ip,
		IpAddress:  mac,
		Rkey:       key,
		Work:       worktype,
		Message:    "",
	}
	data := send_packet.Fluent()
	data = clientsearchsend.AppendByteMsg(data, msg)
	encrypt_buf := make([]byte, len(data))
	C_AES.Encryptbuffer(data, len(data), encrypt_buf)
	_, err := conn.Write(encrypt_buf)
	if err != nil {
		logger.Error(key + ":: Error sending data to server: " + err.Error())
	}
	return nil
}

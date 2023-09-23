package mockagent

import (
	"edetector_go/config"
	"edetector_go/internal/C_AES"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"net"
)

func SendTCPtoServer(worktype task.TaskType, msg string, conn net.Conn) error {
	var send_packet = packet.WorkPacket{
		MacAddress: config.Viper.GetString("MOCK_AGENT_MAC"),
		IpAddress:  config.Viper.GetString("MOCK_AGENT_IP"),
		Rkey:       config.Viper.GetString("MOCK_AGENT_KEY"),
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
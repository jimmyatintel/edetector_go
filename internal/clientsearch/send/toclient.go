package clientsearchsend

import (
	// "context"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/taskchannel"
	"edetector_go/internal/packet"

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

func SendUserTCPtoClient(key string, packet packet.WorkPacket) error{
	task_chan := taskchannel.Task_channel[key]
	task_chan <- packet
	return nil
}
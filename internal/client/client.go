package client

import (
	"edetector_go/internal/client/clientinfo"
	"edetector_go/internal/packet"
)

func PacketClientInfo(p packet.WorkPacket) (clientinfo.ClientInfo, error) {
	clientinfo := clientinfo.ClientInfo{}
	err := clientinfo.Load_data(p.GetMessage())
	return clientinfo, err
}

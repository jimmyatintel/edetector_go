package client

import (
	"edetector_go/internal/client/clientinfo"
	"edetector_go/internal/packet"
)

func init() {

}
func PacketClientInfo(p packet.WorkPacket) clientinfo.ClientInfo {
	clientinfo := clientinfo.ClientInfo{}
	clientinfo.Load_data(p.GetMessage())
	return clientinfo
}

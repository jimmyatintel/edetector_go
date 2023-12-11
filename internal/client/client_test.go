package client

import (
	"edetector_go/internal/client/clientinfo"
	"edetector_go/internal/packet"
	"testing"
)

func TestPacketClientInfo(t *testing.T) {
	test := []struct {
		p    packet.WorkPacket
		want clientinfo.ClientInfo
	}{
		{
			p: packet.WorkPacket{
				Message: "x64|Windows 10 Home|MSI|SYSTEM|1.0.4,1988,1989|800291|3e716e2d61ba910983cb456817116799|0",
			},
			want: clientinfo.ClientInfo{
				SysInfo:      "x64",
				OsInfo:       "Windows 10 Home",
				ComputerName: "MSI",
				UserName:     "SYSTEM",
				FileVersion:  "1.0.4,1988,1989",
				BootTime:     "800291",
				KeyNum:       "3e716e2d61ba910983cb456817116799",
			},
		},
		{
			p: packet.WorkPacket{
				Message: "x64|Windows 10 Home|MSI|SYSTEM|1.0.4,1988,1989",
			},
			want: clientinfo.ClientInfo{
				SysInfo:      "x64",
				OsInfo:       "Windows 10 Home",
				ComputerName: "MSI",
				UserName:     "SYSTEM",
				FileVersion:  "1.0.4,1988,1989",
				BootTime:     "",
				KeyNum:       "",
			},
		},
	}
	for _, tt := range test {
		if got := PacketClientInfo(tt.p); got != tt.want {
			t.Errorf("PacketClientInfo() = %v, want %v", got, tt.want)
		}
	}
}

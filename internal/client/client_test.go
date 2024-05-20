package client

import (
	"edetector_go/internal/client/clientinfo"
	"edetector_go/internal/packet"
	"errors"
	"testing"
)

func TestPacketClientInfo(t *testing.T) {
	test := []struct {
		p    packet.WorkPacket
		want clientinfo.ClientInfo
		err  error
	}{
		{
			p: packet.WorkPacket{
				Message: "x64|@|Windows 10 Home|@|MSI|@|SYSTEM|@|1.0.4,1988,1989|@|800291|@|3e716e2d61ba910983cb456817116799|@|0",
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
			err: nil,
		},
		{
			p: packet.WorkPacket{
				Message: "x64|@|Windows 10 Home|@|MSI|@|SYSTEM|@|1.0.4,1988,1989",
			},
			want: clientinfo.ClientInfo{
				SysInfo:      "",
				OsInfo:       "",
				ComputerName: "",
				UserName:     "",
				FileVersion:  "",
				BootTime:     "",
				KeyNum:       "",
			},
			err: errors.New("error in GiveInfo format, version conflicted"),
		},
	}
	for _, tt := range test {
		got, err := PacketClientInfo(tt.p)
		if err != nil && tt.err == nil {
			t.Errorf("Unexpected error: " + err.Error())
		} else if got != tt.want {
			t.Errorf("PacketClientInfo() = %v, want %v", got, tt.want)
		}
	}
}

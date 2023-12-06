package packet

import (
	"bytes"
	"edetector_go/pkg/file"
	"os"
	"testing"
)

var pktShort []byte
var pktLong []byte
var pktTask []byte

func init() {
	for i := 0; i < 2; i++ {
		file.MoveToParentDir()
	}
	var err error
	pktShort, err = os.ReadFile("test/decrypted_short_packet")
	if err != nil {
		panic(err)
	}
	pktLong, err = os.ReadFile("test/decrypted_long_packet")
	if err != nil {
		panic(err)
	}
	pktTask, err = os.ReadFile("test/decrypted_task_packet")
	if err != nil {
		panic(err)
	}
}

func TestNewWorkPacket(t *testing.T) {
	tests := []struct {
		packet *WorkPacket
		data   []byte
	}{
		{
			data: pktShort,
			packet: &WorkPacket{
				MacAddress: "B8-8A-60-A0-B7-CF",
				IpAddress:  "192.168.200.170",
				Rkey:       "014e4553bbce4a0d85826937225af881",
				Work:       "GiveDetectInfoFirst",
				Message:    "0|0",
			},
		},
	}

	for ind, tt := range tests {
		emptyPkt := &WorkPacket{}
		emptyPkt.NewPacket(tt.data, []byte{})
		if emptyPkt.GetMacAddress() != tt.packet.GetMacAddress() {
			t.Errorf("Failed TestCase %v: GetMacAddress", ind)
		}
		if emptyPkt.GetipAddress() != tt.packet.GetipAddress() {
			t.Errorf("Failed TestCase %v: GetipAddress", ind)
		}
		if emptyPkt.GetRkey() != tt.packet.GetRkey() {
			t.Errorf("Failed TestCase %v: GetRkey", ind)
		}
		if emptyPkt.GetTaskType() != tt.packet.GetTaskType() {
			t.Errorf("Failed TestCase %v: GetTaskType", ind)
		}
		if emptyPkt.GetMessage() != tt.packet.GetMessage() {
			t.Errorf("Failed TestCase %v: GetMessage", ind)
		}
	}
}

func TestNewDataPacket(t *testing.T) {
	tests := []struct {
		packet *DataPacket
		data   []byte
	}{
		{
			data: pktLong,
			packet: &DataPacket{
				MacAddress: "B8-8A-60-A0-B7-CF",
				IpAddress:  "192.168.200.170",
				Rkey:       "014e4553bbce4a0d85826937225af881",
				Work:       "GiveScanProgress",
				Message:    "12/153",
				Raw_data:   pktLong,
			},
		},
	}

	for ind, tt := range tests {
		emptyPkt := &DataPacket{}
		emptyPkt.NewPacket(tt.data, pktLong)
		if emptyPkt.GetMacAddress() != tt.packet.GetMacAddress() {
			t.Errorf("Failed TestCase %v: GetMacAddress", ind)
		}
		if emptyPkt.GetipAddress() != tt.packet.GetipAddress() {
			t.Errorf("Failed TestCase %v: GetipAddress", ind)
		}
		if emptyPkt.GetRkey() != tt.packet.GetRkey() {
			t.Errorf("Failed TestCase %v: GetRkey", ind)
		}
		if emptyPkt.GetTaskType() != tt.packet.GetTaskType() {
			t.Errorf("Failed TestCase %v: GetTaskType", ind)
		}
		if emptyPkt.GetMessage() != tt.packet.GetMessage() {
			t.Errorf("Failed TestCase %v: GetMessage", ind)
		}
		if !bytes.Equal(emptyPkt.Raw_data, tt.packet.Raw_data) {
			t.Errorf("Failed TestCase %v: RawData", ind)
		}
	}
}

func TestNewTaskPacket(t *testing.T) {
	tests := []struct {
		packet *TaskPacket
		data   []byte
	}{
		{
			data: pktTask,
			packet: &TaskPacket{
				Key:     "6b75775ef8854658a595286f6f051399",
				Work:    "StartScan",
				User:    5,
				Message: "Ring0Process",
				Data:    []byte(""),
			},
		},
	}

	for ind, tt := range tests {
		emptyPkt := &TaskPacket{}
		emptyPkt.NewPacket(tt.data)
		if emptyPkt.GetRkey() != tt.packet.GetRkey() {
			t.Errorf("Failed TestCase %v: GetKey", ind)
		}
		if emptyPkt.GetUserTaskType() != tt.packet.GetUserTaskType() {
			t.Errorf("Failed TestCase %v: GetUserTaskType", ind)
		}
		if emptyPkt.User != tt.packet.User {
			t.Errorf("Failed TestCase %v: GetUser", ind)
		}
		if emptyPkt.GetMessage() != tt.packet.GetMessage() {
			t.Errorf("Failed TestCase %v: GetMessage", ind)
		}
		if !bytes.Equal(emptyPkt.Data, tt.packet.Data) {
			t.Errorf("Failed TestCase %v: Data", ind)
		}
	}
}

func TestTaskFluenet(t *testing.T) {
	tests := []struct {
		packet *TaskPacket
		want   []byte
	}{
		{
			want: pktTask,
			packet: &TaskPacket{
				Key:     "6b75775ef8854658a595286f6f051399",
				Work:    "StartScan",
				User:    5,
				Message: "Ring0Process",
				Data:    []byte(""),
			},
		},
	}

	for ind, tt := range tests {
		data := tt.packet.Fluent()
		if !bytes.Equal(data, tt.want) {
			t.Errorf("Failed TestCase %v: Fluent", ind)
		}
	}
}
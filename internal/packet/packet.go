package packet

import (
	"bytes"
	"errors"

	// "github.com/google/uuid"
	"edetector_go/internal/task"
	"edetector_go/pkg/mariadb/query"
	"strings"
)

type Packet interface {
	NewPacket(data []byte, buf []byte) error
	GetMessage() string
	GetMacAddress() string
	GetipAddress() string
	Fluent() []byte
	GetTaskType() task.TaskType
}
type WorkPacket struct {
	// Packet is a struct that contains the packet data
	MacAddress string
	IpAddress  string
	Rkey       string
	Work       task.TaskType
	Message    string
}
type DataPacket struct {
	// Packet is a struct that contains the packet data
	MacAddress string
	IpAddress  string
	Work       task.TaskType
	Message    string
	Raw_data   []byte
}
type TaskPacket struct {
	// Packet is a struct that contains the packet data
	Key     string
	Work    task.TaskType
	User    string
	Message string
	Data    []byte
}

// var Uuid uuid.UUID

func init() {
	// Uuid = uuid.New()
}

func CheckIsWork(p Packet) WorkPacket {
	var NewPacket = p.(*WorkPacket)
	return *NewPacket
}

func CheckIsData(p Packet) DataPacket {
	var NewPacket = p.(*DataPacket)
	return *NewPacket
}

func (p *WorkPacket) NewPacket(data []byte, buf []byte) error {
	error := errors.New("invalid work packet")
	nullIndex := bytes.IndexByte(data[0:20], 0)
	if nullIndex == -1 {
		return error
	}
	p.MacAddress = string(data[0:nullIndex])
	nullIndex = bytes.IndexByte(data[20:40], 0)
	if nullIndex == -1 {
		return error
	}
	p.IpAddress = string(data[20 : 20+nullIndex])
	nullIndex = bytes.IndexByte(data[40:64], 0)
	if nullIndex == -1 {
		return error
	}
	// fmt.Println("function", string(data[40:40+nullIndex]))
	p.Work = task.GetTaskType(string(data[40 : 40+nullIndex]))
	nullIndex = bytes.IndexByte(data[64:], 0)
	// if p.Work == task.UNDEFINE {
	// 	type_error := errors.New("invalid task type" + string(data[40 : 40+nullIndex]))
	// 	return type_error
	// }
	if nullIndex == -1 {
		return error
	}
	p.Message = string(data[64 : 64+nullIndex])
	return nil
}
func (p *DataPacket) NewPacket(data []byte, buf []byte) error {
	error := errors.New("invalid data packet")
	nullIndex := bytes.IndexByte(data[0:20], 0)
	if nullIndex == -1 {
		return error
	}
	p.MacAddress = string(data[0:nullIndex])
	nullIndex = bytes.IndexByte(data[20:40], 0)
	if nullIndex == -1 {
		return error
	}
	p.IpAddress = string(data[20 : 20+nullIndex])
	nullIndex = bytes.IndexByte(data[40:64], 0)
	if nullIndex == -1 {
		return error
	}
	p.Work = task.GetTaskType(string(data[40 : 40+nullIndex]))
	nullIndex = bytes.IndexByte(data[64:], 0)
	if nullIndex == -1 {
		p.Message = string(data[64:])

	} else {
		p.Message = string(data[64 : 64+nullIndex])
	}
	p.Raw_data = buf
	return nil
}
func (p *TaskPacket) NewPacket(data []byte) error {
	error := errors.New("invalid Task packet")
	nullIndex := bytes.IndexByte(data[0:32], 0)
	if nullIndex == -1 {
		return error
	}
	p.Key = string(data[0:nullIndex])
	nullIndex = bytes.IndexByte(data[32:56], 0)
	if nullIndex == -1 {
		return error
	}
	p.Work = task.GetTaskType(string(data[32 : 32+nullIndex]))
	nullIndex = bytes.IndexByte(data[56:76], 0)
	if nullIndex == -1 {
		return error
	}
	p.User = string(data[56 : 56+nullIndex])
	nullIndex = bytes.IndexByte(data[76:], 0)
	if nullIndex == -1 {
		return error
	}
	p.Message = string(data[76 : 76+nullIndex])
	return nil
}
func ljust(s string, width int, fillChar string) []byte {
	data := []byte(s)
	if len(s) >= width {
		return data[0:width]
	}
	data = append(data, 0)
	data = append(data, []byte(strings.Repeat(string(fillChar), width-len(s)-1))...)
	return data
}
func (p *WorkPacket) Fluent() []byte {
	data := []byte("")
	data = append(data, ljust(p.MacAddress, 20, " ")...)
	data = append(data, ljust(p.IpAddress, 20, " ")...)
	data = append(data, ljust(string(p.Work), 24, " ")...)
	data = append(data, ljust(p.Message, 960, " ")...)
	return data
}
func (p *DataPacket) Fluent() []byte {
	data := []byte("")
	data = append(data, ljust(p.MacAddress, 20, " ")...)
	data = append(data, ljust(p.IpAddress, 20, " ")...)
	data = append(data, ljust(string(p.Work), 24, " ")...)
	data = append(data, ljust(p.Message, 960, " ")...)
	return data
}
func (p *TaskPacket) Fluent() []byte {
	data := []byte("")
	data = append(data, ljust(p.Key, 32, " ")...)
	data = append(data, ljust(string(p.Work), 24, " ")...)
	data = append(data, ljust(p.User, 20, " ")...)
	data = append(data, ljust(p.Message, 948, " ")...)
	return data
}
func (p *WorkPacket) GetMessage() string {
	return p.Message
}
func (p *DataPacket) GetMessage() string {
	return p.Message
}
func (p *TaskPacket) GetMessage() string {
	return p.Message
}
func (p *WorkPacket) GetTaskType() task.TaskType {
	return p.Work
}
func (p *DataPacket) GetTaskType() task.TaskType {
	return p.Work
}
func (p *TaskPacket) GetTaskType() task.TaskType {
	return p.Work
}
func (p *WorkPacket) GetMacAddress() string {
	return p.MacAddress
}
func (p *DataPacket) GetMacAddress() string {
	return p.MacAddress
}
func (p *TaskPacket) GetMacAddress() string {
	return query.GetMachineMAC(p.Key)
}
func (p *WorkPacket) GetipAddress() string {
	return p.IpAddress
}
func (p *DataPacket) GetipAddress() string {
	return p.IpAddress
}
func (p *TaskPacket) GetipAddress() string {
	return query.GetMachineIP(p.Key)
}
func (p *WorkPacket) GetRKey() string {
	return p.Rkey
}
func (p *DataPacket) GetRKey() string {
	return "null"
}
func (p *TaskPacket) GetRKey() string {
	return p.Key
}

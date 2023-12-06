package packet

import (
	"bytes"
	"errors"
	"strconv"

	"edetector_go/internal/task"
	"edetector_go/pkg/mariadb/query"
	"encoding/json"
	"net"
	"strings"
)

type Packet interface {
	NewPacket(data []byte, buf []byte) error
	GetMessage() string
	GetMacAddress() string
	GetipAddress() string
	Fluent() []byte
	GetTaskType() task.TaskType
	GetRkey() string
}
type UserPacket interface {
	NewPacket(data []byte) error
	GetMessage() string
	GetMacAddress() string
	GetipAddress() string
	Fluent() []byte
	GetUserTaskType() task.UserTaskType
	GetRkey() string
	Respond(conn net.Conn, isSuccess bool, message string) error
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
	Rkey       string
	Work       task.TaskType
	Message    string
	Raw_data   []byte
}
type TaskPacket struct {
	// Packet is a struct that contains the packet data
	Key     string            `json:"key"`
	Work    task.UserTaskType `json:"work"`
	User    int               `json:"user"`
	Message string            `json:"message"`
	Data    []byte            `json:"data"`
}
type AckPacket struct {
	// Packet is a struct that contains the packet data
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

func CheckIsWork(p Packet) WorkPacket {
	var NewPacket = p.(*WorkPacket)
	return *NewPacket
}

func CheckIsData(p Packet) DataPacket {
	var NewPacket = p.(*DataPacket)
	return *NewPacket
}

// commit
func (p *WorkPacket) NewPacket(data []byte, buf []byte) error {
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
	nullIndex = bytes.IndexByte(data[40:76], 0)
	if nullIndex == -1 {
		return error
	}
	p.Rkey = string(data[40 : 40+nullIndex])
	nullIndex = bytes.IndexByte(data[76:100], 0)
	if nullIndex == -1 {
		return error
	}
	p.Work = task.GetTaskType(string(data[76 : 76+nullIndex]))
	nullIndex = bytes.IndexByte(data[100:], 0)
	if nullIndex == -1 {
		p.Message = string(data[100:])

	} else {
		p.Message = string(data[100 : 100+nullIndex])
	}
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
	nullIndex = bytes.IndexByte(data[40:76], 0)
	if nullIndex == -1 {
		return error
	}
	p.Rkey = string(data[40 : 40+nullIndex])
	nullIndex = bytes.IndexByte(data[76:100], 0)
	if nullIndex == -1 {
		return error
	}
	p.Work = task.GetTaskType(string(data[76 : 76+nullIndex]))
	nullIndex = bytes.IndexByte(data[100:], 0)
	if nullIndex == -1 {
		p.Message = string(data[100:])

	} else {
		p.Message = string(data[100 : 100+nullIndex])
	}
	p.Raw_data = buf
	return nil
}
//To-Do
func (p *TaskPacket) NewPacket(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	p.Work = task.UserTaskType(p.Work)

	// fmt.Println("parse: ", string(p.Key), string(p.Work), string(p.User), string(p.Message))
	return nil
}
func (p *TaskPacket) Respond(conn net.Conn, isSuccess bool, message string) error {
	res := AckPacket{
		IsSuccess: isSuccess,
		Message:   message,
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}
	_, err = conn.Write(resJSON)
	if err != nil {
		return err
	}
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
//To-Do
func (p *WorkPacket) Fluent() []byte {
	data := []byte("")
	data = append(data, ljust(p.MacAddress, 20, " ")...)
	data = append(data, ljust(p.IpAddress, 20, " ")...)
	data = append(data, ljust(p.Rkey, 36, " ")...)
	data = append(data, ljust(string(p.Work), 24, " ")...)
	data = append(data, ljust(p.Message, 924, " ")...)
	return data
}
//To-Do
func (p *DataPacket) Fluent() []byte {
	data := []byte("")
	data = append(data, ljust(p.MacAddress, 20, " ")...)
	data = append(data, ljust(p.IpAddress, 20, " ")...)
	data = append(data, ljust(p.Rkey, 36, " ")...)
	data = append(data, ljust(string(p.Work), 24, " ")...)
	data = append(data, ljust(p.Message, 65436, " ")...)
	// data = append(data, ljust(p.Message, 924, " ")...)
	return data
}
//To-Do
func (p *TaskPacket) Fluent() []byte {
	data := []byte("")
	data = append(data, ljust(p.Key, 32, " ")...)
	data = append(data, ljust(string(p.Work), 24, " ")...)
	data = append(data, ljust(strconv.Itoa(p.User), 20, " ")...)
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
func (p *TaskPacket) GetUserTaskType() task.UserTaskType {
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
func (p *WorkPacket) GetRkey() string {
	return p.Rkey
}
func (p *DataPacket) GetRkey() string {
	return p.Rkey
}
func (p *TaskPacket) GetRkey() string {
	return p.Key
}

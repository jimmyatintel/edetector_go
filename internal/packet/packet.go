package packet

import (
	"bytes"
	"errors"
	"fmt"

	// "github.com/google/uuid"
	"edetector_go/internal/task"
	"edetector_go/pkg/mariadb/query"
	"strings"
	"encoding/json"
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
	Key     string
	Work    task.UserTaskType
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
func (p *TaskPacket) NewPacket(data []byte) error {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}
	key, ok := jsonData["key"].(string)
	if !ok {
		return errors.New("invalid key field")
	}
	p.Key = key

	work, ok := jsonData["work"].(string)
	if !ok {
		return errors.New("invalid work field")
	}
	p.Work = task.UserTaskType(work)

	user, ok := jsonData["user"].(string)
	if !ok {
		return errors.New("invalid user field")
	}
	p.User = user

	messageObj, ok := jsonData["message"].(map[string]interface{})
	if !ok {
		return errors.New("invalid message field")
	}

	process, ok := messageObj["process"].(bool)
	if !ok {
		return errors.New("invalid process field in message")
	}
	network, ok := messageObj["network"].(bool)
	if !ok {
		return errors.New("invalid network field in message")
	}
	p.Message = fmt.Sprintf("%v|%v", process, network)
	fmt.Println("parse: ", string(p.Key), string(p.Work), string(p.User), string(p.Message))

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
	data = append(data, ljust(p.Rkey, 36, " ")...)
	data = append(data, ljust(string(p.Work), 24, " ")...)
	data = append(data, ljust(p.Message, 924, " ")...)
	return data
}
func (p *DataPacket) Fluent() []byte {
	data := []byte("")
	data = append(data, ljust(p.MacAddress, 20, " ")...)
	data = append(data, ljust(p.IpAddress, 20, " ")...)
	data = append(data, ljust(p.Rkey, 36, " ")...)
	data = append(data, ljust(string(p.Work), 24, " ")...)
	data = append(data, ljust(p.Message, 924, " ")...)
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

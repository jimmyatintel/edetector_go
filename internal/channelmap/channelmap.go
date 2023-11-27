package channelmap

import (
	packet "edetector_go/internal/packet"
	"errors"
	"sync"
)

var TaskMu *sync.Mutex
var TaskWorkerChannel map[string](*chan packet.Packet)

var DiskMu *sync.Mutex
var UserDiskChannel = make(map[string](*chan string))

func init() {
	TaskMu = &sync.Mutex{}
	DiskMu = &sync.Mutex{}
}

func AssignTaskChannel(key string, task_chan *chan packet.Packet) {
	TaskMu.Lock()
	TaskWorkerChannel[key] = task_chan
	TaskMu.Unlock()
}

func GetTaskChannel(key string) (chan packet.Packet, error) {
	TaskMu.Lock()
	_, exists := TaskWorkerChannel[key]
	if !exists {
		return nil, errors.New("invalid key")
	}
	task_chan := *TaskWorkerChannel[key]
	TaskMu.Unlock()
	return task_chan, nil
}

func AssignDiskChannel(key string, disk_chan *chan string) {
	DiskMu.Lock()
	UserDiskChannel[key] = disk_chan
	DiskMu.Unlock()
}

func GetDiskChannel(key string) (chan string, error) {
	DiskMu.Lock()
	_, exists := UserDiskChannel[key]
	if !exists {
		return nil, errors.New("invalid key")
	}
	disk_chan := *UserDiskChannel[key]
	DiskMu.Unlock()
	return disk_chan, nil
}

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

var LoadDumpMu *sync.Mutex
var LoadDumpTaskChannel = make(map[string](*chan string))

func init() {
	TaskMu = &sync.Mutex{}
	DiskMu = &sync.Mutex{}
	LoadDumpMu = &sync.Mutex{}
}

func AssignTaskChannel(key string, task_chan *chan packet.Packet) {
	TaskMu.Lock()
	TaskWorkerChannel[key] = task_chan
	TaskMu.Unlock()
}

func GetTaskChannel(key string) (chan packet.Packet, error) {
	TaskMu.Lock()
	_, exists := TaskWorkerChannel[key]
	TaskMu.Unlock()
	if !exists {
		return nil, errors.New("invalid key")
	}
	TaskMu.Lock()
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
	DiskMu.Unlock()
	if !exists {
		return nil, errors.New("invalid key")
	}
	DiskMu.Lock()
	disk_chan := *UserDiskChannel[key]
	DiskMu.Unlock()
	return disk_chan, nil
}

func AssignLoadDumpChannel(key string, dump_chan *chan string) {
	LoadDumpMu.Lock()
	LoadDumpTaskChannel[key] = dump_chan
	LoadDumpMu.Unlock()
}

func GetLoadDumpChannel(key string) (chan string, error) {
	LoadDumpMu.Lock()
	_, exists := LoadDumpTaskChannel[key]
	LoadDumpMu.Unlock()
	if !exists {
		return nil, errors.New("invalid key")
	}

	LoadDumpMu.Lock()
	dump_chan := *LoadDumpTaskChannel[key]
	LoadDumpMu.Unlock()

	return dump_chan, nil
}

func RemoveLoadDumpChannel(key string) error {
	LoadDumpMu.Lock()
	_, exists := LoadDumpTaskChannel[key]
	LoadDumpMu.Unlock()
	if !exists {
		return errors.New("invalid key")
	}

	LoadDumpMu.Lock()
	delete(LoadDumpTaskChannel, key)
	LoadDumpMu.Unlock()

	return nil
}

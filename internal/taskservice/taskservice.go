package taskservice

import (
	"bytes"
	"context"
	"fmt"

	// "edetector_go/internal/fflag"
	"edetector_go/internal/packet"
	work_from_api "edetector_go/internal/work_from_api"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"time"

	"go.uber.org/zap"
)

var ctx context.Context

var taskchans = make(map[string]chan string)

func Start(ctx context.Context) {
	// if enable, err := fflag.FFLAG.FeatureEnabled("taskservice_enable"); enable && err == nil {
		time.Sleep(20 * time.Second)
		fmt.Println("Task service enable.")
		go Main(ctx)
	// } else if err != nil{
	// 	logger.Error("Task service error:", zap.Any("error", err.Error()))
	// } else {
	// 	fmt.Println("Task service disable.")
	// 	return
	// }
}

func Main(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task service is shutting down...")
			return
		default:
			findtask(ctx)
			time.Sleep(5 * time.Second)
		}
	}
}

func findtask(ctx context.Context) {
	unhandle_task := loadallunhandletask()
	for _, task := range unhandle_task {
		content := redis.Redis_get(task.clientid)
		status := redis.GetStatus(content)
		fmt.Println(status)
		if status == 1 {
			if task.clientid == "8beba472f3f44cabbbb44fd232171933" {
				if _, ok := taskchans[task.clientid]; !ok {
					taskchans[task.clientid] = make(chan string, 1024)
					go taskhandler(taskchans[task.clientid], task.clientid, ctx)
				}
				taskchans[task.clientid] <- task.taskid
				Change_task_status(task.taskid, 1)
			}
		}
	}
}

func taskhandler(ch chan string, client string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task handler for " + client + " is shutting down...")
			return
		case taskid := <-ch:
			logger.Info("Task handler for " + client + " is handling task " + taskid)
			message := redis.Redis_get(taskid)
			b := []byte(message)
			Change_task_status(taskid, 2)
			handleTaskrequest(b, taskid, client)
		}
	}
}

func handleTaskrequest(content []byte, taskid string, client string) {
	reqLen := len(content)
	NewPacket := new(packet.TaskPacket)
	err := NewPacket.NewPacket(content)
	if err != nil {
		logger.Error("Error reading task packet:", zap.Any("error", err.Error()), zap.Any("len", reqLen))
		return
	}
	if NewPacket.GetUserTaskType() == "Undefine" {
		nullIndex := bytes.IndexByte(content[76:100], 0)
		logger.Error("Undefine User Task Type: ", zap.String("error", string(content[76:76+nullIndex])))
		return
	}
	logger.Info("Receive task from user", zap.Any("function", NewPacket.GetUserTaskType()))
	taskFunc, ok := work_from_api.WorkapiMap[NewPacket.GetUserTaskType()]
	if !ok {
		logger.Error("Function notfound:", zap.Any("name", NewPacket.GetUserTaskType()))
		return
	}
	_, err = taskFunc(NewPacket)
	if err != nil {
		logger.Error("Task Failed:", zap.Any("error", err.Error()))
		return
	}
	if NewPacket.GetUserTaskType() == "ChangeDetectMode" {
		Change_task_status(taskid, 3)
	}
}

func Finish_task(clientid string, tasktype string) {
	taskID := Find_task_id(clientid, tasktype)
	Change_task_status(taskID, 3)
	Change_task_timestamp(clientid, tasktype)
}

func Stop() {
	ctx.Done()
}

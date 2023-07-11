package taskservice

import (
	"bytes"
	"context"
	"edetector_go/internal/fflag"
	"edetector_go/internal/packet"
	work_from_api "edetector_go/internal/work_from_api"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"time"

	"go.uber.org/zap"
)

var ctx context.Context

var taskchans map[string]chan string

func Start() {
	if enable, err := fflag.FFLAG.FeatureEnabled("taskservice_enable"); enable && err == nil {
		logger.Info("Task service enable.")
		go Main(ctx)
	} else {
		logger.Info("Task service disable.")
		return
	}
}

func Main(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Task service is shutting down...")
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
		if _, ok := taskchans[task.agentid]; !ok {
			taskchans[task.agentid] = make(chan string, 1024)
			go taskhandler(task.agentid, taskchans[task.agentid], ctx)
		}
		taskchans[task.agentid] <- task.taskid
		change_task_status(task.taskid, task.agentid, 1)
	}
}

func taskhandler(agent string, ch chan string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Task handler for " + agent + " is shutting down...")
			return
		case taskid := <-ch:
			logger.Info("Task handler for " + agent + " is handling task " + taskid)
			message := redis.Redis_get(taskid)
			b := []byte(message)
			handleTaskrequest(b)
		}
	}

}

func handleTaskrequest(content []byte) {
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
	_, err = work_from_api.WorkapiMap[NewPacket.GetUserTaskType()](NewPacket)
	if err != nil {
		logger.Error("Function notfound:", zap.Any("name", NewPacket.GetUserTaskType()), zap.Any("error", err.Error()))
		return

	}
}

func Stop() {
	ctx.Done()
}

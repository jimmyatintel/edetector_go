package taskservice

import (
	"context"
	"edetector_go/internal/clientsearch"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"time"
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
			clientsearch.HandleTaskrequest(b)
		}
	}

}
func Stop() {
	ctx.Done()
}

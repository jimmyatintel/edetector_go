package rabbitmq

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	elaInsert "edetector_go/pkg/elastic/insert"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"encoding/json"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Index string `json:"index"`
	Data  string `json:"data"`
}

func ToRabbitMQ_Details(index string, st elastic.Request_data, values []string, uuid string, agentID string, ip string, name string, item string, date string, ttype string, etc string, priority string, taskType string) error {
	taskID := query.Load_task_id(agentID, taskType, 2)
	template, err := StringToStruct(st, values, uuid, agentID, ip, name, item, date, ttype, etc, taskID)
	if err != nil {
		return err
	}
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	var msg = Message{
		Index: index,
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for {
		err = Publish(priority, msgBytes)
		if err != nil {
			logger.Error("Error sending to rabbitMQ (details), retrying... " + err.Error())
			randomSleep := (rand.Intn(100) + 1) * 100 // 0.1 ~ 10
			time.Sleep(time.Duration(randomSleep) * time.Millisecond)
		} else {
			break
		}
	}
	return nil
}

func ToRabbitMQ_Relation(index string, template elastic.Request_data, priority string) error {
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	var msg = Message{
		Index: config.Viper.GetString("ELASTIC_PREFIX") + index,
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for {
		err = Publish(priority, msgBytes)
		if err != nil {
			logger.Error("Error sending to rabbitMQ (relation), retrying... " + err.Error())
			randomSleep := (rand.Intn(100) + 1) * 100 // 0.1 ~ 10
			time.Sleep(time.Duration(randomSleep) * time.Millisecond)
		} else {
			break
		}
	}
	return nil
}

func ToRabbitMQ_FinishSignal(agent string, taskType string, priority string) error {
	logger.Info("Finish signal sent to rabbitMQ")
	template := elaInsert.FinishSignal{
		Agent:    agent,
		TaskType: taskType,
	}
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	var msg = Message{
		Index: "Finish",
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for {
		err = Publish(priority, msgBytes)
		if err != nil {
			logger.Error("Error sending to rabbitMQ (main), retrying... " + err.Error())
			randomSleep := (rand.Intn(100) + 1) * 100 // 0.1 ~ 10
			time.Sleep(time.Duration(randomSleep) * time.Millisecond)
		} else {
			break
		}
	}
	return nil
}

func StringToStruct(st elastic.Request_data, values []string, uuid string, agentID string, ip string, name string, item string, date string, ttype string, etc string, taskID string) (elastic.Request_data, error) {
	v := reflect.Indirect(reflect.ValueOf(st))
	values = append(values, uuid, agentID, ip, name, item, date, ttype, etc, taskID)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Int:
			values[i] = strings.TrimSpace(values[i])
			value, err := strconv.Atoi(values[i])
			if err != nil {
				logger.Error("Error converting to int: " + err.Error())
				break
			}
			field.Set(reflect.ValueOf(value))
		case reflect.Int64:
			values[i] = strings.TrimSpace(values[i])
			value, err := strconv.ParseInt(values[i], 10, 64)
			if err != nil {
				logger.Error("Error converting to int64: " + err.Error())
				break
			}
			field.Set(reflect.ValueOf(value))
		case reflect.String:
			field.Set(reflect.ValueOf(values[i]))
		case reflect.Bool:
			value, err := strconv.ParseBool(values[i])
			if err != nil {
				logger.Error("Error converting to bool: " + err.Error())
				break
			}
			field.Set(reflect.ValueOf(value))
		}
	}
	return st, nil
}

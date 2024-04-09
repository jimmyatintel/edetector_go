package rabbitmq

import (
	"edetector_go/pkg/elastic"
	elaInsert "edetector_go/pkg/elastic/insert"
	"edetector_go/pkg/logger"
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

func ToRabbitMQ_Details(index string, st elastic.Request_data, sub_st elastic.Request_data, values []string, uuid string, agentID string, ip string, name string, item string, date string, ttype string, etc string, priority string, taskType string, taskID string, category string) error {
	values = append(values, uuid, agentID, ip, name, item, date, ttype, etc, taskID, category)
	template, err := StringToStruct(st, sub_st, values)
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

func ToRabbitMQ_Tree(index string, template elastic.Request_data, priority string) error {
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

func StringToStruct(st elastic.Request_data, sub_st elastic.Request_data, values []string) (elastic.Request_data, error) {
	v := reflect.Indirect(reflect.ValueOf(st))
	j := 0
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Struct:
			nestedValues := make([]string, 0)
			nestedV := reflect.Indirect(reflect.ValueOf(sub_st))
			for j = i; j < nestedV.NumField()+i; j++ {
				nestedValues = append(nestedValues, values[j])
			}
			j--
			// Create new instance of nested struct
			nestedSt, err := StringToStruct(sub_st, nil, nestedValues)
			if err != nil {
				return nil, err
			}
			field.Set(reflect.ValueOf(nestedSt).Elem())
		case reflect.Int:
			values[j] = strings.TrimSpace(values[j])
			value, err := strconv.Atoi(values[j])
			if err != nil {
				logger.Error("Error converting to int [" + strconv.Itoa(j) + "]: " + err.Error())
			}
			field.Set(reflect.ValueOf(value))
		case reflect.Int64:
			values[j] = strings.TrimSpace(values[j])
			value, err := strconv.ParseInt(values[j], 10, 64)
			if err != nil {
				logger.Error("Error converting to int64 [" + strconv.Itoa(j) + "]: " + err.Error())
			}
			field.Set(reflect.ValueOf(value))
		case reflect.String:
			field.Set(reflect.ValueOf(values[j]))
		case reflect.Bool:
			value, err := strconv.ParseBool(values[j])
			if err != nil {
				logger.Error("Error converting to bool [" + strconv.Itoa(j) + "]: " + err.Error())
			}
			field.Set(reflect.ValueOf(value))
		}
		j++
	}
	return st, nil
}

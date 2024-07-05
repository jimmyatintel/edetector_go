package insert

import (
	"edetector_go/pkg/elastic"
	elaDelete "edetector_go/pkg/elastic/delete"
	"edetector_go/pkg/logger"
	mariadbquery "edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"encoding/json"
	"strings"
)

type IndexInfo struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
}

type FinishSignal struct {
	Agent    string `json:"agent"`
	TaskType string `json:"taskType"`
	TaskID   string `json:"taskID"`
}

func (s *FinishSignal) Elastical() ([]byte, error) {
	return json.Marshal(s)
}

func BulkInsert(action []string, work []string) error {
	var buf strings.Builder
	var finishSignals []string
	for i, doc := range action {
		if IsFinish(doc) {
			finishSignals = append(finishSignals, work[i])
			continue
		}
		buf.WriteString(doc)
		buf.WriteByte('\n')
		buf.WriteString(work[i])
		buf.WriteByte('\n')
	}
	if buf.Len() != 0 {
		err := elastic.BulkIndexRequest(buf, 0)
		if err != nil {
			return err
		}
	}
	for _, signal := range finishSignals {
		var data FinishSignal
		err := json.Unmarshal([]byte(signal), &data)
		if err != nil {
			logger.Error("Error unmarshaling finish signal: " + err.Error())
			continue
		}
		logger.Info("Finish signal received: " + data.Agent + " " + data.TaskType)
		task_id := mariadbquery.Load_task_id(data.Agent, data.TaskType, 2)
		if task_id != data.TaskID {
			logger.Warn("Task ID mismatch: " + task_id + " " + data.TaskID)
			continue
		}
		if data.TaskType == "StartGetDrive" || data.TaskType == "StartMemoryTree" { // delete head first
			logger.Debug("Delete old TreeHead " + data.TaskType + ": " + data.Agent)
			err = elaDelete.DeleteOldData(data.Agent, data.TaskType, task_id, true)
			if err != nil {
				logger.Error("Error deleting TreeHead " + data.TaskType + ": " + err.Error())
			}
		}
		if data.TaskType == "StartGetDrive" || data.TaskType == "StartMemoryTree" || data.TaskType == "StartCollect" {
			logger.Debug("Delete old repeated data: " + data.Agent + " " + data.TaskType)
			err = elaDelete.DeleteOldData(data.Agent, data.TaskType, task_id, false)
			if err != nil {
				logger.Error("Error deleting old repeated data: " + err.Error())
			}
		}
		mariadbquery.Finish_task(data.Agent, data.TaskType)

		// delete task info in redis after finish
		if err := redis.RedisDelete(data.TaskID); err != nil {
			logger.Error("Error deleting task info in redis: " + err.Error())
		} else {
			logger.Info("Task info deleted in redis: " + data.TaskID)
		}
	}
	return nil
}

func IsFinish(jsonString string) bool {
	var indexInfo struct {
		IndexInfo `json:"index"`
	}
	err := json.Unmarshal([]byte(jsonString), &indexInfo)
	if err != nil {
		logger.Error("Error parsing action: " + err.Error())
		return false
	}
	if indexInfo.Index == "Finish" {
		return true
	}
	return false
}

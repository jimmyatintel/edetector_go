package insert

import (
	"edetector_go/pkg/elastic"
	elaDelete "edetector_go/pkg/elastic/delete"
	"edetector_go/pkg/logger"
	mariadbquery "edetector_go/pkg/mariadb/query"
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
}

func (s *FinishSignal) Elastical() ([]byte, error) {
	return json.Marshal(s)
}

func BulkInsert(action []string, work []string) error {
	var buf strings.Builder
	for i, doc := range action {
		if IsFinish(doc) {
			var data FinishSignal
			err := json.Unmarshal([]byte(work[i]), &data)
			if err != nil {
				logger.Error("Error unmarshaling finish signal: " + err.Error())
				continue
			}
			logger.Info("Finish signal received: " + data.Agent + " " + data.TaskType)
			task_id := mariadbquery.Load_task_id(data.Agent, data.TaskType, 2)
			if data.TaskType == "StartGetDrive" { // delete head first
				logger.Info("Delete old ExplorerTreeHead")
				err = elaDelete.DeleteOldData(data.Agent, "ExplorerTreeHead", task_id)
				if err != nil {
					logger.Error("Error deleting ExplorerTreeHead: " + err.Error())
				}
			}
			if data.TaskType != "StartScan" {
				logger.Info("Delete old repeated data" + data.Agent + " " + data.TaskType)
				err = elaDelete.DeleteOldData(data.Agent, data.TaskType, task_id)
				if err != nil {
					logger.Error("Error deleting old repeated data: " + err.Error())
				}
			}
			mariadbquery.Finish_task(data.Agent, data.TaskType)
			continue
		}
		buf.WriteString(doc)
		buf.WriteByte('\n')
		buf.WriteString(work[i])
		buf.WriteByte('\n')
	}
	if buf.Len() == 0 {
		return nil
	}
	err := elastic.BulkIndexRequest(buf)
	if err != nil {
		return err
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

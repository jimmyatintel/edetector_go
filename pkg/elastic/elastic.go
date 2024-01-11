package elastic

import (
	"context"
	"edetector_go/config"
	"edetector_go/pkg/logger"
	mariadbquery "edetector_go/pkg/mariadb/query"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"go.uber.org/zap"
)

var es *elasticsearch.Client

type IndexInfo struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
}

func flagcheck() bool {
	// if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
	return true
	// }
	// return false
}

func Elastic_init() {
	var err error
	cfg := elasticsearch.Config{
		Addresses: []string{"http://" + config.Viper.GetString("ELASTIC_HOST") + ":" + config.Viper.GetString("ELASTIC_PORT")},
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Panic("Error connecting to elastic: " + err.Error())
		panic(err)
	}
}

func CreateIndex(name string) {
	if !flagcheck() {
		return
	}
	req := esapi.IndicesCreateRequest{
		Index: name,
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error creating index: " + err.Error())
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		logger.Error("Error creating index: " + res.String())
		return
	}
	logger.Info("Created index: " + res.String())
}

func IndexRequest(name string, body string) error {
	if !flagcheck() {
		return nil
	}
	req := esapi.IndexRequest{
		Index: name,
		Body:  strings.NewReader(body),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return errors.New("Error indexing: " + res.String())
	}
	logger.Debug("Index request: " + res.String())
	return nil
}

func BulkIndexRequest(action []string, work []string) error {
	if !flagcheck() {
		return nil
	}
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
	res, err := es.Bulk(
		strings.NewReader(buf.String()),
		es.Bulk.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return errors.New("Error BulkIndexRequest: " + res.String())
	}
	output := res.String()
	if output[:8] == "[200 OK]" {
		logger.Info("BulkIndexRequest Res: " + output[:100])
	} else {
		return errors.New("Error BulkIndexRequest Res: " + output[:500])
	}
	return nil
}

func UpdateByQueryRequest(query string, index string) (int, error) {
	if !flagcheck() {
		return 0, nil
	}
	updateReq := esapi.UpdateByQueryRequest{
		Index: []string{index},
		Body:  strings.NewReader(query),
	}
	updateRes, err := updateReq.Do(context.Background(), es)
	if err != nil {
		return 0, err
	}
	defer updateRes.Body.Close()
	if updateRes.IsError() {
		return 0, errors.New("Error UpdateByQueryRequest:" + updateRes.String())
	}
	var updateResponse map[string]interface{}
	if err := json.NewDecoder(updateRes.Body).Decode(&updateResponse); err != nil {
		return 0, err
	}
	updated, found := updateResponse["updated"]
	if !found {
		return 0, fmt.Errorf("updated count not found in the response")
	}
	updatedFloat, ok := updated.(float64)
	if !ok {
		return 0, fmt.Errorf("updated count is not a number")
	}
	return int(updatedFloat), nil
}

func UpdateByDocIDRequest(index string, docID string, script string) error {
	if !flagcheck() {
		return nil
	}
	updateReq := esapi.UpdateRequest{
		Index:      index,
		DocumentID: docID,
		Body:       strings.NewReader(script),
	}
	updateRes, err := updateReq.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer updateRes.Body.Close()
	if updateRes.IsError() {
		return errors.New("Error UpdateByDocIDRequest:" + updateRes.String())
	}
	logger.Debug("Update response" + updateRes.String())
	return nil
}

func SearchRequest(index string, body string, sortItem string) []interface{} {
	if !flagcheck() {
		return nil
	}
	from := 0
	size := 10
	hitsArray := []interface{}{}
	var result map[string]interface{}
	for {
		req := esapi.SearchRequest{
			Index: []string{index},
			Body:  strings.NewReader(body),
			From:  &from,
			Size:  &size,
			Sort:  []string{sortItem},
		}
		res, err := req.Do(context.Background(), es)
		if err != nil {
			logger.Error("Error getting response: " + err.Error())
			return nil
		}
		defer res.Body.Close()
		if res.IsError() {
			logger.Error("Error SearchRequest: " + res.String())
			return nil
		}
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			logger.Error("Error decoding response: " + err.Error())
			return nil
		}
		hits, ok := result["hits"].(map[string]interface{})
		if !ok {
			logger.Error("Hits not found in response: " + res.String())
			return nil
		}
		newHitsArray, ok := hits["hits"].([]interface{})
		if !ok {
			logger.Error("Hits array not found in response")
			return nil
		}
		hitsArray = append(hitsArray, newHitsArray...)
		if len(newHitsArray) == 0 || sortItem == "" {
			return hitsArray
		}
		from += 10
	}
}

func DeleteByQueryRequest(indexes []string, deleteQuery string) error {
	if !flagcheck() {
		return errors.New("elastic is not enabled")
	}
	req := esapi.DeleteByQueryRequest{
		Index: indexes,
		Body:  strings.NewReader(deleteQuery),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return errors.New("Error DeleteByQueryRequest:" + res.String())
	} else {
		var responseJSON map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&responseJSON)
		if err != nil {
			return err
		}
		logger.Info("Deleted repeated data: ", zap.Any("message", responseJSON["deleted"]))

		conflictCount := responseJSON["version_conflicts"].(float64)
		if conflictCount != 0 {
			logger.Error("Version conflict: ", zap.Any("message", conflictCount))
		}
		failures := responseJSON["failures"].([]interface{})
		if len(failures) != 0 {
			logger.Error("Failures: ", zap.Any("message", failures))
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

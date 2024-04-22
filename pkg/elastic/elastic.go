package elastic

import (
	"context"
	"edetector_go/config"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"go.uber.org/zap"
)

var es *elasticsearch.Client

func elasticCheck(retry int) bool {
	if retry > 2 {
		logger.Error("Retry count exceeded")
		return false
	}
	for {
		if redis.RedisGetInt("elastic_restart") == 1 {
			logger.Warn("Elastic is restarting, waiting for 60 seconds")
			time.Sleep(60 * time.Second)
		} else {
			return true
		}
	}
	// if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
	// return true
	// }
	// return false
}

func Elastic_init() {
	var err error
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Viper.GetString("ELASTIC_URL_1"),
			config.Viper.GetString("ELASTIC_URL_2"),
			config.Viper.GetString("ELASTIC_URL_3"),
		},
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Panic("Error connecting to elastic: " + err.Error())
		panic(err)
	}
}

type Request_data interface {
	Elastical() ([]byte, error)
}

type BulkResponse struct {
	Took   int         `json:"took"`
	Errors bool        `json:"errors"`
	Items  []ItemIndex `json:"items"`
}

type ItemIndex struct {
	Index IndexDetails `json:"index"`
}

type IndexDetails struct {
	Index  string `json:"_index"`
	Type   string `json:"_type"`
	ID     string `json:"_id"`
	Status int    `json:"status"`
	Error  Error  `json:"error,omitempty"`
}

type Error struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

func CreateIndex(name string, retry int) {
	if !elasticCheck(retry) {
		logger.Error("elastic is not enabled or exceed retry count")
		return
	}
	req := esapi.IndicesCreateRequest{
		Index: name,
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error creating index: " + err.Error())
		time.Sleep(60 * time.Second)
		CreateIndex(name, retry+1)
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		logger.Error("Error creating index: " + res.String())
		return
	}
	logger.Info("Created index: " + res.String())
}

func IndexRequest(name string, body string, retry int) error {
	if !elasticCheck(retry) {
		return errors.New("elastic is not enabled or exceed retry count")
	}
	req := esapi.IndexRequest{
		Index: name,
		Body:  strings.NewReader(body),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error indexing: " + err.Error())
		time.Sleep(60 * time.Second)
		return IndexRequest(name, body, retry+1)
	}
	defer res.Body.Close()
	if res.IsError() {
		return errors.New("Error indexing: " + res.String())
	}
	logger.Debug("Index request: " + res.String())
	return nil
}

func BulkIndexRequest(buf strings.Builder, retry int) error {
	if !elasticCheck(retry) {
		return errors.New("elastic is not enabled or exceed retry count")
	}
	res, err := es.Bulk(
		strings.NewReader(buf.String()),
		es.Bulk.WithContext(context.Background()),
	)
	if err != nil {
		logger.Error("Error BulkIndexRequest: " + err.Error())
		time.Sleep(60 * time.Second)
		return BulkIndexRequest(buf, retry+1)
	}
	defer res.Body.Close()
	if res.IsError() {
		return errors.New("Error BulkIndexRequest: " + res.String())
	}
	output := res.String()
	if output[:8] == "[200 OK]" {
		logger.Info("BulkIndexRequest Res: " + output[:100])
		var bulkResp BulkResponse
		err := json.NewDecoder(res.Body).Decode(&bulkResp)
		if err != nil {
			logger.Error("Error unmarshalling BulkResponse: " + err.Error())
		} else if bulkResp.Errors {
			for _, item := range bulkResp.Items {
				if item.Index.Status != 201 {
					logger.Error("Error indexing: " + item.Index.Error.Reason)
				}
			}
		}
	} else {
		return errors.New("Error BulkIndexRequest Res: " + output[:500])
	}
	return nil
}

func UpdateByQueryRequest(query string, index string, retry int) (int, error) {
	if !elasticCheck(retry) {
		return 0, errors.New("elastic is not enabled or exceed retry count")
	}
	updateReq := esapi.UpdateByQueryRequest{
		Index: []string{index},
		Body:  strings.NewReader(query),
	}
	updateRes, err := updateReq.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error UpdateByQueryRequest: " + err.Error())
		time.Sleep(60 * time.Second)
		return UpdateByQueryRequest(query, index, retry+1)
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

func UpdateByDocIDRequest(index string, docID string, script string, retry int) error {
	if !elasticCheck(retry) {
		return errors.New("elastic is not enabled or exceed retry count")
	}
	updateReq := esapi.UpdateRequest{
		Index:      index,
		DocumentID: docID,
		Body:       strings.NewReader(script),
	}
	updateRes, err := updateReq.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error UpdateByDocIDRequest: " + err.Error())
		time.Sleep(60 * time.Second)
		return UpdateByDocIDRequest(index, docID, script, retry+1)
	}
	defer updateRes.Body.Close()
	if updateRes.IsError() {
		return errors.New("Error UpdateByDocIDRequest:" + updateRes.String())
	}
	logger.Debug("Update query: " + script)
	logger.Debug("Update response: " + updateRes.String())
	return nil
}

func SearchRequest(index string, body string, sortItem string, retry int) []interface{} {
	if !elasticCheck(retry) {
		return nil
	}
	var result map[string]interface{}
	size := 10000
	req := esapi.SearchRequest{
		Index:  []string{index},
		Body:   strings.NewReader(body),
		Scroll: 60 * time.Second,
		Size:   &size,
		Sort:   []string{sortItem},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error getting response: " + err.Error())
		time.Sleep(60 * time.Second)
		return SearchRequest(index, body, sortItem, retry+1)
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
	scrollID, ok := result["_scroll_id"].(string)
	if !ok {
		logger.Error("ScrollID not found in response: " + res.String())
		return nil
	}
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		logger.Error("Hits not found in response: " + res.String())
		return nil
	}
	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		logger.Error("Hits array not found in response")
		return nil
	}
	for {
		req := esapi.ScrollRequest{
			Scroll:   60 * time.Second,
			ScrollID: scrollID,
		}
		res, err := req.Do(context.Background(), es)
		if err != nil {
			logger.Error("Error getting response: " + err.Error())
			return hitsArray
		}
		defer res.Body.Close()
		if res.IsError() {
			logger.Error("Error ScrollRequest: " + res.String())
			return hitsArray
		}
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			logger.Error("Error decoding response: " + err.Error())
			return hitsArray
		}
		hits, ok = result["hits"].(map[string]interface{})
		if !ok {
			logger.Error("Hits not found in response: " + res.String())
			return hitsArray
		}
		newHitsArray, ok := hits["hits"].([]interface{})
		if !ok {
			logger.Error("Hits array not found in response")
			return hitsArray
		}
		if len(newHitsArray) == 0 {
			break
		}
		hitsArray = append(hitsArray, newHitsArray...)
	}
	return hitsArray
}

func DeleteByQueryRequest(indexes []string, query string, retry int) error {
	if !elasticCheck(retry) {
		return errors.New("elastic is not enabled")
	}
	logger.Debug("Index: " + strings.Join(indexes, ", "))
	logger.Debug("Delete query:" + query)
	req := esapi.DeleteByQueryRequest{
		Index: indexes,
		Body:  strings.NewReader(query),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error DeleteByQueryRequest: " + err.Error())
		time.Sleep(60 * time.Second)
		return DeleteByQueryRequest(indexes, query, retry+1)
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
		logger.Info("Delete By Query: ", zap.Any("message", responseJSON["deleted"]))

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

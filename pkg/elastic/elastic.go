package elastic

import (
	"context"
	"edetector_go/config"
	"edetector_go/pkg/fflag"
	"edetector_go/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"go.uber.org/zap"
)

var es *elasticsearch.Client

var diskIndex = []string{"explorer", "explorer_relation"}

var dbIndex = []string{"AppResourceUsageMonitor", "ARPCache", "BaseService", "ChromeBookmarks", "ChromeCache", "ChromeDownload",
	"ChromeHistory", "ChromeKeywordSearch", "ChromeLogin", "DNSInfo", "EdgeBookmarks", "EdgeCache", "EdgeCookies", "EdgeHistory",
	"EdgeLogin", "EventApplication", "EventSecurity", "EventSystem", "FirefoxBookmarks", "FirefoxCache", "FirefoxCookies",
	"FirefoxHistory", "IEHistory", "InstalledSoftware", "JumpList", "MUICache", "Network", "NetworkDataUsageMonitor",
	"NetworkResources", "OpenedFiles", "Prefetch", "Process", "Service", "Shortcuts", "StartRun", "TaskSchedule",
	"USBdevices", "UserAssist", "UserProfiles", "WindowsActivity"}

func flagcheck() bool {
	if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
		return true
	}
	return false
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
		logger.Error("error creating index: " + err.Error())
	}
	defer res.Body.Close()
	logger.Info(res.String())
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
	logger.Debug("Index content: ", zap.Any("message", body))
	logger.Debug("Index request: ", zap.Any("message", res.String()))
	return nil
}

func BulkIndexRequest(action []string, work []string) error {
	if !flagcheck() {
		return nil
	}
	var buf strings.Builder
	for i, doc := range action {
		buf.WriteString(doc)
		buf.WriteByte('\n')
		buf.WriteString(work[i])
		buf.WriteByte('\n')
	}
	res, err := es.Bulk(
		strings.NewReader(buf.String()),
		es.Bulk.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	logger.Info("len:", zap.Any("message", len(action)))
	index := 0
	for {
		ind := strings.Index(res.String()[index:], "error")
		if ind == -1 {
			break
		}
		output := ""
		if ind+300 > len(res.String()) {
			output = res.String()[index:]
		} else {
			output = res.String()[index : index+300]
		}
		logger.Info("res: ", zap.Any("message", output))
		index = ind + 1
	}
	return nil
}

func UpdateRequest(agent string, id string, time string, index string) error {
	if !flagcheck() {
		return nil
	}
	query := fmt.Sprintf(`
	{
		"script": {
			"source": "ctx._source.processConnectIP = params.processConnectIP",
			"lang": "painless",
			"params": {
				"processConnectIP": "true"
			}
		},
		"query": {
			"bool": {
				"must": [
					{ "term": { "agent": "%s" } },
					{ "term": { "processId": %s } },
					{ "term": { "processCreateTime": %s } },
					{ "term": { "mode": "detect" } }
				]
			}
		}
	}`, agent, id, time)
	updateReq := esapi.UpdateByQueryRequest{
		Index: []string{index},
		Body:  strings.NewReader(query),
	}
	updateRes, err := updateReq.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer updateRes.Body.Close()
	if updateRes.IsError() {
		responseBytes, _ := ioutil.ReadAll(updateRes.Body)
		errorMessage := string(responseBytes)
		return errors.New(errorMessage)
	} else {
		var response map[string]interface{}
		if err := json.NewDecoder(updateRes.Body).Decode(&response); err != nil {
			return err
		}
		updatedCount := int(response["updated"].(float64))
		if updatedCount > 0 {
			logger.Debug("Update network of the detect process: ", zap.Any("message", agent+" "+id+" "+time))
		} else {
			createBody := fmt.Sprintf(`
			{
				"agent": "%s",
				"processId": %s,
				"processCreateTime": %s,
				"processConnectIP": "true",
				"mode": "detect"
			}`, agent, id, time)
			err = IndexRequest(index, createBody)
			if err != nil {
				return err
			}
			logger.Debug("Create a new detect process: ", zap.Any("message", agent+" "+id+" "+time))
		}
	}
	return nil
}

func SearchRequest(index string, body string) string {
	if !flagcheck() {
		return ""
	}
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(body),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	// logger.Info(res.String())
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logger.Error("Error decoding response: ", zap.Any("error", err.Error()))
		return ""
	}
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		logger.Error("Hits not found in response")
		return ""
	}
	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		logger.Error("Hits array not found in response")
		return ""
	}
	var docID string
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			logger.Error("Hit is not a map")
			continue
		}
		docIDVal, ok := hitMap["_id"].(string)
		if !ok {
			logger.Error("Doc ID not found in hit")
			continue
		}
		docID = docIDVal
		break
	}
	return docID
}

func DeleteByQueryRequest(field string, value string, ttype string) error {
	if !flagcheck() {
		return errors.New("elastic is not enabled")
	}
	deleteQuery := fmt.Sprintf(`
	{
		"query": {
			"term": {
				"%s": "%s"
			}
		}
	}
	`, field, value)
	req := esapi.DeleteByQueryRequest{
		Index: getIndexes(ttype),
		Body:  strings.NewReader(deleteQuery),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("error response")
	} else {
		var responseJSON map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&responseJSON)
		if err != nil {
			return err
		}
		logger.Info("Deleted repeated data:", zap.Any("message", responseJSON["deleted"]))

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

func getIndexes(ttype string) []string {
	prefix := config.Viper.GetString("ELASTIC_PREFIX")
	indexes := []string{}
	switch ttype {
	case "StartGetDrive":
		for _, ind := range diskIndex {
			indexes = append(indexes, prefix+"_"+ind)
		}
	case "StartCollect":
		for _, ind := range dbIndex {
			indexes = append(indexes, prefix+"_"+strings.ToLower(ind))
		}
	}
	return indexes
}

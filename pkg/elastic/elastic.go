package elastic

import (
	"context"
	"edetector_go/config"
	"edetector_go/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
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
	"USBdevices", "UserAssist", "UserProfiles", "WindowsActivity", "Wireless"}

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
	}
	defer res.Body.Close()
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
	logger.Debug("Index request: " + res.String())
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
	output := res.String()
	if output[:8] == "[200 OK]" {
		logger.Info("BulkIndexRequest Res: " + output[:100])
	} else {
		return errors.New("Error BulkIndexRequest Res: " + output[:100])
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
		return 0, errors.New(updateRes.String())
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

func UpdateByDocIDRequest(index string, docID string, newValue string, source string) error {
	if !flagcheck() {
		return nil
	}
	script := map[string]interface{}{
		"script": map[string]interface{}{
			"source": source,
			"lang":   "painless",
			"params": map[string]interface{}{
				"value": newValue,
			},
		},
	}
	scriptBytes, err := json.Marshal(script)
	if err != nil {
		return err
	}
	updateReq := esapi.UpdateRequest{
		Index:      index,
		DocumentID: docID,
		Body:       strings.NewReader(string(scriptBytes)),
	}

	updateRes, err := updateReq.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer updateRes.Body.Close()
	return nil
}

func SearchRequest(index string, body string) []interface{} {
	if !flagcheck() {
		return nil
	}
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(body),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error getting response: " + err.Error())
		return nil
	}
	defer res.Body.Close()
	// logger.Info(res.String())
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logger.Error("Error decoding response: " + err.Error())
		return nil
	}
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		logger.Error("Hits not found in response")
		return nil
	}
	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		logger.Error("Hits array not found in response")
		return nil
	}
	return hitsArray
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
		logger.Info("Deleted repeated data ("+field+"-"+value+"): ", zap.Any("message", responseJSON["deleted"]))

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

package elastic

import (
	"bytes"
	"context"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/logger"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"go.uber.org/zap"
)

var es *elasticsearch.Client

func flagcheck() bool {
	if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
		return true
	}
	return false
}
func SetElkClient() error {
	var err error
	cfg := elasticsearch.Config{
		Addresses: []string{"http://" + os.Getenv("ELASTIC_HOST") + ":" + os.Getenv("ELASTIC_PORT")},
	}
	es, err = elasticsearch.NewClient(cfg)
	return err
}

func createIndex(name string) {
	if !flagcheck() {
		return
	}
	req := esapi.IndicesCreateRequest{
		Index: name,
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		logger.Error(err.Error())
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

func BulkUpdateDocuments(index string, docIDs []string) {
	var buf bytes.Buffer
	for _, docID := range docIDs {
		action := map[string]interface{}{
			"update": map[string]interface{}{
				"_index": index,
				"_id":    docID,
				"_type":  "_doc",
			},
		}
		source, err := json.Marshal(action)
		if err != nil {
			logger.Info("Failed to marshal", zap.Any("message", err.Error()))
		}
		buf.Write(source)
		buf.WriteByte('\n')
		updateData := map[string]interface{}{
			"doc": map[string]interface{}{
				"processConnectIP": "detected",
			},
		}
		docUpdateData, err := json.Marshal(updateData)
		if err != nil {
			logger.Info("Failed to marshal", zap.Any("message", err.Error()))
		}
		buf.Write(docUpdateData)
		buf.WriteByte('\n')
	}
	res, err := es.Bulk(
		strings.NewReader(buf.String()),
		es.Bulk.WithContext(context.Background()),
	)
	if err != nil {
		logger.Info("Failed to marshal", zap.Any("message", err.Error()))
	}
	defer res.Body.Close()
	logger.Info("Bulk update: ", zap.Any("message", res.String()))
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
	var docID string
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		docID = hit.(map[string]interface{})["_id"].(string)
	}
	return docID
}

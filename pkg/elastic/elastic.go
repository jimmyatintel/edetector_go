package elastic

import (
	"context"
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/logger"
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
		Addresses: []string{"http://" + config.Viper.GetString("ELASTIC_HOST") + ":" + config.Viper.GetString("ELASTIC_PORT")},
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
	// logger.Info("Bulk Index request: ", zap.Any("message", action))
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

func updateRequest(index string, body string) {
	if !flagcheck() {
		return
	}
	updateReq := esapi.UpdateRequest{
		Index: index,
		Body:  strings.NewReader(body),
	}
	res, err := updateReq.Do(context.Background(), es)
	if err != nil {
		logger.Error("Error executing update request: %s", zap.Any("error", err.Error()))
	}
	defer res.Body.Close()
	if res.IsError() {
		logger.Error("Error response: ", zap.Any("message", res.Status()))
	}
	logger.Info(res.String())
}

func searchRequest(index string, body string) {
	if !flagcheck() {
		return
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
	logger.Info(res.String())
}

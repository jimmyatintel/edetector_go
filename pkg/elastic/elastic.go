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
	ind := strings.Index(res.String(), "error")
	if ind != -1 {
		output := ""
		if ind + 1 > len(res.String()) {
			output = res.String()[ind:]
		} else {
			output = res.String()[ind:ind+300]
		}
		logger.Info("error: ", zap.Any("message", output))
	} else {
		logger.Info("sucess")
	}
	return nil
}

func searchRequest(name string, body string) {
	if !flagcheck() {
		return
	}
	req := esapi.SearchRequest{
		Index: []string{name},
		Body:  strings.NewReader(body),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	logger.Info(res.String())
}

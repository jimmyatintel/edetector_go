package elasticquery

import (
	elastic "edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

func UpdateNetworkInfo(agent string, networkSet map[string]struct{}) {
outerLoop:
	for {
		var docs []string
		for key := range networkSet {
			values := strings.Split(key, ",")
			pid := values[0]
			createTime := values[1]
			query := fmt.Sprintf(`{
				"query": {
					"bool": {
						"must": [
							{ "term": { "agent": "%s" } },
							{ "term": { "processId": "%s" } },
							{ "term": { "processCreateTime": "%s" } }
						]
					}
				}
			}`, agent, pid, createTime)
			logger.Info("network info", zap.Any("message", agent+" "+pid+" "+createTime))
			doc := elastic.SearchRequest("ed_de_memory", query)
			docs = append(docs, doc)
			if doc == "" {
				logger.Info("waiting for 10 seconds for process info")
				time.Sleep(10 * time.Second)
				continue outerLoop
			}
		}
		logger.Info("docs", zap.Any("docs", docs))
		elastic.BulkUpdateDocuments("ed_de_memory", docs)
	}
}

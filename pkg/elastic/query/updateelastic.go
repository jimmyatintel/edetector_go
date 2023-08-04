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
fail_count := 0
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
			doc := elastic.SearchRequest("ed_memory", query)
			docs = append(docs, doc)
			if doc == "" {
				logger.Info("waiting 60s for updating process: ", zap.Any("message", pid+" "+createTime))
				time.Sleep(60 * time.Second)
				fail_count += 1
				if fail_count >= 3 {
                    logger.Error("fail to update process: ", zap.Any("message", pid+" "+createTime))
                    break outerLoop
                }
				continue outerLoop
			}
		}
		elastic.BulkUpdateDocuments("ed_memory", docs)
		break
	}
}

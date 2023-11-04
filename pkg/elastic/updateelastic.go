package elastic

import (
	"edetector_go/config"
	"edetector_go/pkg/logger"
	"fmt"
	"strings"
)

func UpdateNetworkInfo(agent string, networkSet map[string]struct{}) {
	for key := range networkSet {
		values := strings.Split(key, ",")
		id := values[0]
		time := values[1]
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
		updatedCount, err := UpdateByQueryRequest(query, config.Viper.GetString("ELASTIC_PREFIX")+"_memory")
		if err != nil {
			logger.Error("Error updating detect process: " + err.Error())
			continue
		}
		if updatedCount > 0 {
			logger.Debug("Update network of the detect process: " + agent + "|" + id + "|" + time)
		} else {
			createBody := fmt.Sprintf(`
			{
				"agent": "%s",
				"processId": %s,
				"processCreateTime": %s,
				"processConnectIP": "true",
				"mode": "detect"
			}`, agent, id, time)
			err = IndexRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", createBody)
			if err != nil {
				logger.Error("Error creating detect process: " + err.Error())
				continue
			}
			logger.Debug("Create a new detect process: " + agent + "|" + id + "|" + time)
		}
	}
}

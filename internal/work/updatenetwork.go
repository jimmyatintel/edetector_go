package work

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"fmt"
	"strconv"
	"strings"
)

func UpdateNetworkInfo(agent string, networkSet map[string]struct{}) {
	for key := range networkSet {
		values := strings.Split(key, ",")
		id := values[0]
		time := values[1]
		logger.Info("UpdateNetworkInfo: " + id + "|" + time)
		malicious, err := strconv.Atoi(values[2])
		if err != nil {
			logger.Error("Error converting malicious to int: " + err.Error())
			continue
		}
		// search for the detect process
		index := config.Viper.GetString("ELASTIC_PREFIX") + "_memory"
		searchQuery := fmt.Sprintf(`{
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
		hitsArray := elastic.SearchRequest(index, searchQuery)
		if len(hitsArray) == 0 { // detect process not exists
			risklevel := scoretoLevel(malicious * 20)
			createBody := fmt.Sprintf(`
			{
				"agent": "%s",
				"processId": %s,
				"processCreateTime": %s,
				"processConnectIP": "true",
				"riskLevel": %d,
				"riskScore": %d,
				"mode": "OnlyNetwork"
			}`, agent, id, time, risklevel, malicious*20)
			err := elastic.IndexRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", createBody)
			if err != nil {
				logger.Error("Error creating detect process: " + err.Error())
				continue
			}
			logger.Debug("Create a new detect process: " + agent + "|" + id + "|" + time)
		} else if len(hitsArray) > 0 { // update the detect process
			if len(hitsArray) > 1 {
				logger.Warn("More than one detect process: " + agent + "|" + id + "|" + time)
			}
			for _, hit := range hitsArray {
				hitMap, ok := hit.(map[string]interface{})
				if !ok {
					logger.Error("Hit is not a map")
					continue
				}
				docID, ok := hitMap["_id"].(string)
				if !ok {
					logger.Error("docID not found")
					continue
				}
				// get the risklevel of the detect process
				source, ok := hitMap["_source"].(map[string]interface{})
				if !ok {
					logger.Error("source not found")
					continue
				}
				score, ok := source["riskScore"].(float64)
				if !ok {
					logger.Error("riskScore not found")
					continue
				}
				riskscore := int(score) + (malicious * 20)
				risklevel := scoretoLevel(riskscore)
				script := fmt.Sprintf(`
				{
					"script": {
						"source": "ctx._source.riskLevel = params.level; ctx._source.riskScore = params.score",
						"lang": "painless",
						"params": {
							"level": %d,
							"score": %d
						}
					}
				}`, risklevel, riskscore)
				err := elastic.UpdateByDocIDRequest(index, docID, script)
				if err != nil {
					logger.Error("Error updating detect process: " + err.Error())
					continue
				}
			}
		}
	}
}

package work

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"errors"
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
		searchDetectQuery := fmt.Sprintf(`{
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
		hitsDetectArray := elastic.SearchRequest(index, searchDetectQuery, "uuid")
		searchNetworkQuery := fmt.Sprintf(`{
			"query": {
				"bool": {
				  "must": [
					{ "term": { "agent": "%s" } },
					{ "term": { "processId": %s } },
					{ "term": { "processCreateTime": %s } },
					{ "term": { "mode": "detectNetwork" } }
				  ]
				}
			  }
			}`, agent, id, time)
		hitsNetworktArray := elastic.SearchRequest(index, searchNetworkQuery, "uuid")
		previousScore := 0.0
		if len(hitsNetworktArray) > 0 {
			previousScore, _, err = getScore(hitsNetworktArray[0])
			if err != nil {
				logger.Error("Error getting score: " + err.Error())
			}
			elastic.DeleteByQueryRequest(searchNetworkQuery, "Memory")
		}
		if len(hitsDetectArray) == 0 { // detect process not exists
			riskscore := int(previousScore) + malicious*20
			risklevel := scoretoLevel(riskscore)
			ip, name, err := query.GetMachineIPandName(agent)
			if err != nil {
				logger.Error("Error getting agent ip and name: " + err.Error())
				continue
			}
			createBody := fmt.Sprintf(`
			{
				"agent": "%s",
				"processId": %s,
				"processCreateTime": %s,
				"processConnectIP": "true",
				"riskLevel": %d,
				"riskScore": %d,
				"processName": "Unknown",
				"agentIP": "%s",
				"agentName": "%s",
				"mode": "detectNetwork"
			}`, agent, id, time, risklevel, riskscore, ip, name)
			err = elastic.IndexRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", createBody)
			if err != nil {
				logger.Error("Error creating detect process: " + err.Error())
				continue
			}
			logger.Debug("Create a new detect process: " + agent + "|" + id + "|" + time)
		} else if len(hitsDetectArray) > 0 { // update the detect process
			if len(hitsDetectArray) > 1 {
				logger.Warn("More than one detect process: " + agent + "|" + id + "|" + time)
			}
			for _, hit := range hitsDetectArray {
				score, docID, error := getScore(hit)
				if error != nil {
					logger.Error("Error getting score: " + error.Error())
					continue
				}
				riskscore := int(previousScore) + int(score) + (malicious * 20)
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

func getScore(hit interface{}) (float64, string, error) {
	hitMap, ok := hit.(map[string]interface{})
	if !ok {
		return 0, "", errors.New("hit is not a map")
	}
	docID, ok := hitMap["_id"].(string)
	if !ok {
		return 0, "", errors.New("docID not found")
	}
	// get the risklevel of the detect process
	source, ok := hitMap["_source"].(map[string]interface{})
	if !ok {
		return 0, "", errors.New("source not found")
	}
	score, ok := source["riskScore"].(float64)
	if !ok {
		return 0, "", errors.New("riskScore not found")
	}
	return score, docID, nil
}

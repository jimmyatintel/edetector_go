package work

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"fmt"
	"strings"
	"time"
)

func HandleRelation(lines []string, key string, count int) {
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\r", "")
		values := strings.Split(line, "|")
		if len(values) != count {
			if len(values) != 1 {
				logger.Error("Invalid line: " + line)
			}
			continue
		}
		BuildMemoryRelation(key, "parent", values[5], values[5], values[9])
		BuildMemoryRelation(key, "child", values[9], values[5], values[9])
		time.Sleep(1 * time.Second)
	}
}

func BuildMemoryRelation(agent string, field string, value string, parent string, child string) {
	index := config.Viper.GetString("ELASTIC_PREFIX") + "_memory_relation"
	searchQuery := fmt.Sprintf(`{
			"query": {
				"bool": {
				  "must": [
					{ "term": { "agent": "%s" } },
					{ "term": { "parent": "%s" } }
				  ]
				}
			  }
			}`, agent, value)
	hitsArray := elastic.SearchRequest(index, searchQuery, "uuid")
	if len(hitsArray) == 0 { // not exists
		data := MemoryRelation{}
		if field == "parent" {
			data = MemoryRelation{
				Agent:  agent,
				IsRoot: true,
				Parent: value,
				Child:  []string{child},
			}
		} else if field == "child" {
			data = MemoryRelation{
				Agent:  agent,
				IsRoot: false,
				Parent: value,
				Child:  []string{},
			}
		}
		err := rabbitmq.ToRabbitMQ_Relation("_memory_relation", data, "ed_high")
		if err != nil {
			logger.Error("Error sending to rabbitMQ (relation): " + err.Error())
			return
		}
	} else if len(hitsArray) == 1 { // exists
		hitMap, ok := hitsArray[0].(map[string]interface{})
		if !ok {
			logger.Error("Hit is not a map")
			return
		}
		docID, ok := hitMap["_id"].(string)
		if !ok {
			logger.Error("docID not found")
			return
		}
		if field == "parent" {
			source, ok := hitMap["_source"].(map[string]interface{})
			if !ok {
				logger.Error("source not found")
				return
			}
			oldChild, ok := source["child"].([]interface{})
			if !ok {
				logger.Error("children not found")
			}
			for _, c := range oldChild {
				if c.(string) == child {
					return
				}
			}
			script := fmt.Sprintf(`
			{
				"script": {
					"source": "ctx._source.child.add(params.value)",
					"lang": "painless",
					"params": {
						"value": "%s"
					}
				}
			}`, child)
			err := elastic.UpdateByDocIDRequest(index, docID, script)
			if err != nil {
				logger.Error("Error updating parent: " + err.Error())
				return
			}
		} else if field == "child" {
			script := `
			{
				"script": {
					"source": "ctx._source.isRoot = params.value",
					"lang": "painless",
					"params": {
						"value": "false"
					}
				}
			}`
			err := elastic.UpdateByDocIDRequest(index, docID, script)
			if err != nil {
				logger.Error("Error updating child: " + err.Error())
				return
			}
		}
	} else {
		logger.Error("More than one relation found: " + field + " parent " + parent + " child " + child)
	}
}

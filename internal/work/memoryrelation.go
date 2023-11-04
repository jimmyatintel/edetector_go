package work

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"fmt"
)

func BuildMemoryRelation(agent string, field string, value string, parent string, child string) {
	index := config.Viper.GetString("ELASTIC_PREFIX") + "_memory_relation"
	searchQuery := fmt.Sprintf(`{
			"query": {
				"bool": {
				  "must": [
					{ "query_string": { "fields": ["agent"], "query": "%s" } },
					{ "query_string": { "fields": ["%s"], "query": "%s" } }
				  ]
				}
			  }
			}`, agent, field, value)
	hitsArray := elastic.SearchRequest(index, searchQuery)
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
			_, err := elastic.UpdateByDocIDRequest(index, docID, child, "ctx._source.child.add(params.value)")
			if err != nil {
				logger.Error("Error updating parent: " + err.Error())
			}
		} else if field == "child" {
			_, err := elastic.UpdateByDocIDRequest(index, docID, child, "ctx._source.isRoot = params.value")
			if err != nil {
				logger.Error("Error updating child: " + err.Error())
			}
		}

	} else {
		logger.Error("More than one relation found" + field + " " + parent + " " + child)
	}
}

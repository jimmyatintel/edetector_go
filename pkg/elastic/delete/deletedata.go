package delete

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"fmt"
)

var diskIndex = []string{"explorer", "explorer_relation"}

func GetIndexes(ttype string) []string {
	prefix := config.Viper.GetString("ELASTIC_PREFIX")
	indexes := []string{}
	switch ttype {
	case "StartGetDriveHead":
		indexes = append(indexes, prefix+"_explorer_relation")
	case "StartMemoryTreeHead":
		indexes = append(indexes, prefix+"_memory_relation")
	case "StartMemoryTree":
		indexes = append(indexes, prefix+"_memory_tree")
		indexes = append(indexes, prefix+"_memory_relation")
	case "StartGetDrive":
		for _, ind := range diskIndex {
			indexes = append(indexes, prefix+"_"+ind)
		}
	case "StartCollect":
		indexes = append(indexes, prefix+"_collection")
	case "Memory":
		indexes = append(indexes, prefix+"_memory")
	}
	return indexes
}

func DeleteOldData(key string, ttype string, taskID string) error {
	indexes := GetIndexes(ttype)
	var query string
	if ttype == "StartGetDriveHead" || ttype == "StartMemoryTreeHead" {
		query = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } },
						{ "term": { "isRoot": true } }
					],
					"must_not": [
						{ "term": { "task_id": "%s" } }
					]
				}
			}
		}`, key, taskID)
	} else {
		query = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } }
					],
					"must_not": [
						{ "term": { "task_id": "%s" } }
					]
				}
			}
		}`, key, taskID)
	}
	err := elastic.DeleteByQueryRequest(indexes, query)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUnfinishedData(key string, ttype string, taskID string) error {
	indexes := GetIndexes(ttype)
	var query string
	if ttype == "ExplorerTreeHead" {
		query = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } },
						{ "term": { "task_id": "%s" } },
						{ "term": { "isRoot": true } }
					]
				}
			}
		}`, key, taskID)
	} else {
		query = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } },
						{ "term": { "task_id": "%s" } }
					]
				}
			}
		}`, key, taskID)
	}
	err := elastic.DeleteByQueryRequest(indexes, query)
	if err != nil {
		return err
	}
	return nil
}

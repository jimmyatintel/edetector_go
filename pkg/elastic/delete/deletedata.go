package delete

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"fmt"
)

func GetIndexes(ttype string) []string {
	prefix := config.Viper.GetString("ELASTIC_PREFIX")
	indexes := []string{}
	switch ttype {
	case "StartGetDrive":
		indexes = append(indexes, prefix+"_explorer")
	case "StartCollect":
		indexes = append(indexes, prefix+"_collection")
	case "StartMemoryTree":
		indexes = append(indexes, prefix+"_memory")
	default:
		return nil
	}
	return indexes
}

func DeleteOldData(key string, ttype string, taskID string, head bool) error {
	indexes := GetIndexes(ttype)
	if indexes == nil {
		return fmt.Errorf("invalid task type")
	}
	var query string
	if head {
		category := "nil"
		if ttype == "StartGetDrive" {
			category = "explorer"
		} else if ttype == "StartMemoryTree" {
			category = "memory_tree"
		}
		query = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } },
						{ "term": { "%s.isRoot": true } }
					],
					"must_not": [
						{ "term": { "task_id": "%s" } }
					]
				}
			}
		}`, key, category, taskID)
	} else if ttype == "StartMemoryTree" {
		query = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } },
						{ "term": { "category": "memory_tree" } }
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
	err := elastic.DeleteByQueryRequest(indexes, query, 0)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUnfinishedData(key string, ttype string, taskID string, head bool) error {
	indexes := GetIndexes(ttype)
	if indexes == nil {
		return fmt.Errorf("invalid task type")
	}
	var query string
	if head {
		category := "nil"
		if ttype == "StartGetDrive" {
			category = "explorer"
		} else if ttype == "StartMemoryTree" {
			category = "memory_tree"
		}
		query = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } },
						{ "term": { "task_id": "%s" } },
						{ "term": { "%s.isRoot": true } }
					]
				}
			}
		}`, key, taskID, category)
	} else if ttype == "StartMemoryTree" {
		query = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } },
						{ "term": { "task_id": "%s" } },
						{ "term": { "category": "memory_tree" } }
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
	err := elastic.DeleteByQueryRequest(indexes, query, 0)
	if err != nil {
		return err
	}
	return nil
}

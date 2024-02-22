package delete

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"fmt"
	"strings"
)

var diskIndex = []string{"explorer", "explorer_relation"}

var dbIndex = []string{"AppResourceUsageMonitor", "ARPCache", "BaseService", "ChromeBookmarks", "ChromeCache", "ChromeDownload",
	"ChromeHistory", "ChromeKeywordSearch", "ChromeLogin", "DNSInfo", "EdgeBookmarks", "EdgeCache", "EdgeCookies", "EdgeHistory",
	"EdgeLogin", "EventApplication", "EventSecurity", "EventSystem", "FirefoxBookmarks", "FirefoxCache", "FirefoxCookies",
	"FirefoxHistory", "IEHistory", "InstalledSoftware", "JumpList", "MUICache", "Network", "NetworkDataUsageMonitor",
	"NetworkResources", "OpenedFiles", "Prefetch", "Process", "Service", "Shortcuts", "StartRun", "TaskSchedule",
	"USBdevices", "UserAssist", "UserProfiles", "WindowsActivity", "Wireless"}

func GetIndexes(ttype string) []string {
	prefix := config.Viper.GetString("ELASTIC_PREFIX")
	indexes := []string{}
	switch ttype {
	case "ExplorerTreeHead":
		indexes = append(indexes, prefix+"_explorer_relation")
	case "StartGetDrive":
		for _, ind := range diskIndex {
			indexes = append(indexes, prefix+"_"+ind)
		}
	case "StartCollect":
		for _, ind := range dbIndex {
			indexes = append(indexes, prefix+"_"+strings.ToLower(ind))
		}
	case "Memory":
		indexes = append(indexes, prefix+"_memory")
	}
	return indexes
}

func DeleteOldData(key string, ttype string, taskID string) error {
	indexes := GetIndexes(ttype)
	var query string
	if ttype == "ExplorerTreeHead" {
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

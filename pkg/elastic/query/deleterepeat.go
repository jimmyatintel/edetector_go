package query

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

func DeleteRepeat(key string, ttype string) {
	indexes := GetIndexes(ttype)
	// detail
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"agent": "%s"
			}
		}
	}`, key)
	elastic.DeleteByQueryRequest(indexes, query)
	// main
	for _, ind := range indexes {
		query = fmt.Sprintf(`{
			"query": {
				"term": { "agent": "%s"},
				"term": { "index": "%s"}
			}
		}`, key, ind)
		mainInd := []string{config.Viper.GetString("ELASTIC_PREFIX") + "_main"}
		elastic.DeleteByQueryRequest(mainInd, query)
	}
}

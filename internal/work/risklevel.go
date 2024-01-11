package work

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"edetector_go/pkg/virustotal"
	"fmt"
	"strconv"
	"strings"
)

type hack_list struct {
	ID          int    `json:"id"`
	ProcessName string `json:"processName"`
	Cmd         string `json:"cmd"`
	Path        string `json:"path"`
	AddingPoint int    `json:"addingPoint"`
}

type white_black_list struct {
	ID        int    `json:"id"`
	FileName  string `json:"filename"`
	MD5       string `json:"md5"`
	Sign      string `json:"signature"`
	Path      string `json:"path"`
	SetupUser string `json:"setupUser"`
	Reason    string `json:"reason"`
}

type HackReq struct {
	Added   []hack_list `json:"added"`
	Removed []hack_list `json:"removed"`
}

type WhiteBlackReq struct {
	Added   []white_black_list `json:"added"`
	Removed []white_black_list `json:"removed"`
}

func Getriskscore(info Memory, initScore int) (string, string, string, string, error) {
	score := initScore
	// virus total
	maliciluos, total, err := virustotal.ScanFile(info.ProcessMD5)
	if err != nil {
		logger.Warn("Error getting virustotal: " + err.Error())
	} else {
		if maliciluos > 0 {
			score += maliciluos * 20
		}
	}
	realPath := strings.Replace(info.ProcessPath, "\\\\", "\\", -1)
	// white list
	whiteList, err := query.Load_white_list()
	if err != nil {
		logger.Error("Error loading white list" + err.Error())
	} else {
		for _, white := range whiteList {
			if (white[0] == "" || info.ProcessName == white[0]) && (white[1] == "" || info.ProcessMD5 == white[1]) && strings.Contains(info.DigitalSign, white[2]) && strings.Contains(realPath, white[3]) {
				logger.Debug("Hit white list")
				return "0", "0", strconv.Itoa(maliciluos), strconv.Itoa(total), nil
			}
		}
	}
	// black list
	blackList, err := query.Load_black_list()
	if err != nil {
		logger.Error("Error loading black list" + err.Error())
	} else {
		for _, black := range blackList {
			if (black[0] == "" || info.ProcessName == black[0]) && (black[1] == "" || info.ProcessMD5 == black[1]) && strings.Contains(info.DigitalSign, black[2]) && strings.Contains(realPath, black[3]) {
				logger.Debug("Hit black list")
				return "3", "300", strconv.Itoa(maliciluos), strconv.Itoa(total), nil
			}
		}
	}
	// hack list
	hackList, err := query.Load_hack_list()
	if err != nil {
		logger.Error("Error loading hack list" + err.Error())
	} else {
		for _, hack := range hackList {
			if (hack[0] == "" || info.ProcessName == hack[0]) && strings.Contains(info.DynamicCommand, hack[1]) && strings.Contains(realPath, hack[2]) {
				point, err := strconv.Atoi(hack[3])
				if err != nil {
					logger.Error("Error converting adding_point to integer" + err.Error())
					continue
				}
				logger.Debug("Hit hack list")
				score += point
			}
		}
	}

	if info.ProcessBeInjected == 2 {
		if _, ok := HighRiskMap[info.ProcessName]; ok {
			score += 150
		} else {
			score += 90
		}
	}
	if info.ProcessBeInjected == 1 {
		score += 30
	}
	if len(info.InjectActive) > 0 && info.InjectActive[0] == '1' && info.DigitalSign == "null" {
		score += 60
	}
	if len(info.InjectActive) > 2 && info.InjectActive[2] == '1' && info.DigitalSign == "null" {
		score += 30
	}
	if len(info.Boot) > 0 && info.Boot[0] == '1' && info.DigitalSign == "null" {
		score += 30
	}
	if len(info.Boot) > 2 && info.Boot[2] == '1' && info.DigitalSign == "null" {
		score += 30
	}
	if info.ProcessConnectIP != "false" {
		score += 30
	}
	if info.ImportOtherDLL != "null" {
		score += 60
	}
	if len(info.Hide) > 0 && info.Hide[0] == '1' {
		score += 150
	}
	if len(info.Hide) > 2 && info.Hide[2] == '1' {
		score += 60
	}
	if strings.Contains(info.Hook, "NtQuerySystemInformation") || strings.Contains(info.Hook, "RtlGetNativeSystemInformation") || strings.Contains(info.Hook, "ZwQuerySystemInformation") {
		score += 150
	}
	if info.ParentProcessPath == "null" {
		if score > 60 {
			score -= 60
		} else {
			score = 0
		}
	}
	if info.DigitalSign == "null" && info.ProcessBeInjected == 0 && info.Hook == "null" {
		score = 0
	}
	level := scoretoLevel(score)
	return strconv.Itoa(level), strconv.Itoa(score), strconv.Itoa(maliciluos), strconv.Itoa(total), nil
}

func scoretoLevel(score int) int {
	if score >= 150 {
		return 3
	} else if score > 90 {
		return 2
	} else if score > 30 {
		return 1
	} else {
		return 0
	}
}

func UpdateLists(listType string, req interface{}) {
	switch listType {
	case "hack":
		hackReq := req.(HackReq)
		for _, add := range hackReq.Added {
			query := buildHackQuery(add.ProcessName, add.Cmd, add.Path)
			recalculateScore(query)
		}
		for _, removed := range hackReq.Removed {
			query := buildHackQuery(removed.ProcessName, removed.Cmd, removed.Path)
			recalculateScore(query)
		}
	case "white":
		fallthrough
	case "black":
		WhiteBlackReq := req.(WhiteBlackReq)
		for _, add := range WhiteBlackReq.Added {
			query := buildWhiteBlackQuery(add.FileName, add.MD5, add.Sign, add.Path)
			recalculateScore(query)
		}
		for _, removed := range WhiteBlackReq.Removed {
			query := buildWhiteBlackQuery(removed.FileName, removed.MD5, removed.Sign, removed.Path)
			recalculateScore(query)
		}
	}
	redis.RedisSet("update_"+listType+"_list", "1")
	logger.Info("Update " + listType + " list finished")
}

func buildHackQuery(name, cmd, path string) string {
	path = strings.ReplaceAll(path, "\\", "\\\\")
	path = strings.ReplaceAll(path, "/", "//")
	return fmt.Sprintf(`{
		"size": 10,
		"query": {
			"bool": {
				"must": [
				{
					"query_string": {
					"fields": ["processName"],
					"query": "%s"
					}
				},
				{
					"query_string": {
					"fields": ["dynamicCommand"],
					"query": "*%s*"
					}
				},
				{
					"query_string": {
					"fields": ["processPath"],
					"query": "*%s*"
					}
				}
				]
			}
		},
		"sort": [
			{
				"uuid": {
					"order": "asc"
				}
			}
		]
	}`, name, cmd, path)
}

func buildWhiteBlackQuery(name, md5, sign, path string) string {
	path = strings.ReplaceAll(path, "\\", "\\\\")
	path = strings.ReplaceAll(path, "/", "//")
	return fmt.Sprintf(`{
		"size": 10,
		"query": {
			"bool": {
				"must": [
				  {
					"query_string": {
					  "fields": ["processName"],
					  "query": "%s"
					}
				  },
				  {
					"query_string": {
					  "fields": ["processMD5"],
					  "query": "*%s*"
					}
				  },
				  {
					"query_string": {
					  "fields": ["digitalSign"],
					  "query": "*%s*"
					}
				  },
				  {
					"query_string": {
						"fields": ["processPath"],
						"query": "*%s*"
					}
				  }
				]
			  }
			}
		  },
		  "sort": [
			  {
				  "uuid": {
					  "order": "asc"
				  }
			  }
		  ]
	  }`, name, md5, sign, path)
}

func recalculateScore(query string) {
	logger.Debug("query: " + query)
	hitsArray := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", query)
	logger.Debug("Hits len: " + strconv.Itoa(len(hitsArray)))
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			logger.Error("Error converting hit to map")
			return
		}
		// convert hitMap to Memory struct
		var info Memory
		info.Mode = hitMap["_source"].(map[string]interface{})["mode"].(string)
		if info.Mode != "scan" && info.Mode != "detect" {
			return
		}
		info.ProcessName = hitMap["_source"].(map[string]interface{})["processName"].(string)
		info.ProcessCreateTime = int(hitMap["_source"].(map[string]interface{})["processCreateTime"].(float64))
		info.DynamicCommand = hitMap["_source"].(map[string]interface{})["dynamicCommand"].(string)
		info.ProcessMD5 = hitMap["_source"].(map[string]interface{})["processMD5"].(string)
		info.ProcessPath = hitMap["_source"].(map[string]interface{})["processPath"].(string)
		info.ParentProcessId = int(hitMap["_source"].(map[string]interface{})["parentProcessId"].(float64))
		info.ParentProcessName = hitMap["_source"].(map[string]interface{})["parentProcessName"].(string)
		info.ParentProcessPath = hitMap["_source"].(map[string]interface{})["parentProcessPath"].(string)
		info.DigitalSign = hitMap["_source"].(map[string]interface{})["digitalSign"].(string)
		info.ProcessId = int(hitMap["_source"].(map[string]interface{})["processId"].(float64))
		info.InjectActive = hitMap["_source"].(map[string]interface{})["injectActive"].(string)
		info.ProcessBeInjected = int(hitMap["_source"].(map[string]interface{})["processBeInjected"].(float64))
		info.Boot = hitMap["_source"].(map[string]interface{})["boot"].(string)
		info.Hide = hitMap["_source"].(map[string]interface{})["hide"].(string)
		info.ImportOtherDLL = hitMap["_source"].(map[string]interface{})["importOtherDLL"].(string)
		info.Hook = hitMap["_source"].(map[string]interface{})["hook"].(string)
		info.ProcessConnectIP = hitMap["_source"].(map[string]interface{})["processConnectIP"].(string)
		info.Agent = hitMap["_source"].(map[string]interface{})["agent"].(string)

		// update the score
		initScore := getNetworkMalicious(info.Agent, info.ProcessId, info.ProcessCreateTime)
		level, score, _, _, err := Getriskscore(info, initScore)
		if err != nil {
			logger.Error("Error getting risk score: " + err.Error())
			return
		}
		docID := hitMap["_id"].(string)
		if !ok {
			logger.Error("docID not found")
			return
		}
		query := fmt.Sprintf(`{
				"script": {
					"source": "ctx._source.riskLevel = params.level; ctx._source.riskScore = params.score",
					"lang": "painless",
					"params": {
						"level": %s,
						"score": %s
					}
				}
			}`, level, score)
		err = elastic.UpdateByDocIDRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", docID, query)
	}

}

func getNetworkMalicious(agent string, pid int, ctime int) int {
	total := 0
	query := fmt.Sprintf(`{
		"query": {
			"bool": {
				"must": [
					{ "term": { "agent": "%s" } },
					{ "term": { "processId": %s } },
					{ "term": { "processCreateTime": %s } }
				]
			}
		}
	}`, agent, strconv.Itoa(pid), strconv.Itoa(ctime))
	hitsArray := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory_network", query)
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			logger.Error("Error converting hit to map")
			return 0
		}
		score := int(hitMap["_source"].(map[string]interface{})["malicious"].(float64))
		if score > 0 {
			total += score * 20
		}
	}
	return total
}

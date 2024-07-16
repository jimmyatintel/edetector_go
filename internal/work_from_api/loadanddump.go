package workfromapi

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/task"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"edetector_go/pkg/request"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var dumpWorkingPath = filepath.Join("static", "dumpWorking")

func StartLoadDll(key, taskId, msg string) (task.TaskResult, error) {
	logger.Info(key + "::StartLoadDll: " + msg)

	// store taskId in redis
	redisKey := key + string(task.START_LOAD_DLL) + strings.Split(msg, "|")[0]
	if err := redis.RedisSet(redisKey, taskId); err != nil {
		logger.Error("StartLoadDll: redis set error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: redisKey,
			LoadDll:  true,
			Progress: -1,
		})
		return task.FAIL, err
	}

	time.Sleep(5 * time.Second)

	err := clientsearchsend.SendUserTCPtoClientUsingKey(key, task.GET_LOAD_DLL, msg)
	if err != nil {
		logger.Error("StartLoadDll: send tcp error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: redisKey,
			LoadDll:  true,
			Progress: -1,
		})
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func StartDumpDll(key, taskId, msg string) (task.TaskResult, error) {
	logger.Info(key + "::StartDumpDll: " + msg)

	// store taskId in redis
	msgs := strings.Split(msg, "|")
	redisKey := key + string(task.START_DUMP_DLL) + msgs[0] + "|" + msgs[1]
	if err := redis.RedisSet(redisKey, taskId); err != nil {
		logger.Error("StartDumpDll: redis set error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: redisKey,
			Progress: -1,
		})
		return task.FAIL, err
	}

	err := clientsearchsend.SendUserTCPtoClientUsingKey(key, task.GET_DUMP_DLL, msg)
	if err != nil {
		logger.Error("StartDumpDll: send tcp error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: redisKey,
			Progress: -1,
		})
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func StartDumpProcess(key, taskId, msg string) (task.TaskResult, error) {
	logger.Info(key + "::StartDumpProcess: " + msg)

	// store taskId in redis
	redisKey := key + string(task.START_DUMP_PROCESS) + strings.Split(msg, "|")[0]
	if err := redis.RedisSet(redisKey, taskId); err != nil {
		logger.Error("StartDumpProcess: redis set error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: redisKey,
			Progress: -1,
		})
		return task.FAIL, err
	}

	err := clientsearchsend.SendUserTCPtoClientUsingKey(key, task.GET_DUMP_PROCESS, msg)
	if err != nil {
		logger.Error("StartDumpProcess: send tcp error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: redisKey,
			Progress: -1,
		})
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func StartDumpDrive(key, taskId, msg string) (task.TaskResult, error) {
	logger.Info(key + "::StartDumpDrive: " + msg)
	msgs := strings.Split(msg, "|")
	fileId, filePath := msgs[0], msgs[1]
	hitPath := strings.TrimRight(filePath, "\r")
	hitPath = strings.ReplaceAll(hitPath, "\\\\", "_backslash_")
	hitPath = strings.ReplaceAll(hitPath, "\\", "_backslash_")
	hitPath = strings.ReplaceAll(hitPath, "//", "_slash_")
	hitPath = strings.ReplaceAll(hitPath, "/", "_slash_")
	hitPath = strings.ReplaceAll(hitPath, ":", "_colon_")

	var followingPath string
	if strings.Contains(hitPath, "_backslash_") {
		followingPath = hitPath + "_backslash_*"
	} else {
		followingPath = hitPath + "_slash_*"
	}
	redisKey := key + string(task.START_DUMP_DRIVE) + filePath

	// declare an error handler to avoid duplicate code
	errHandler := func(deleteFile bool, err string) {
		if deleteFile {
			os.Remove(filepath.Join(dumpWorkingPath, taskId+".txt"))
		}
		logger.Error("StartDumpDrive: "+err, logger.GetCallerInfoForLog()...)
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			RedisKey: redisKey,
			Progress: -1,
		})
	}

	// store taskId in redis
	if err := redis.RedisSet(redisKey, taskId); err != nil {
		errHandler(false, "redis error: "+err.Error())
		return task.FAIL, err
	}

	// open a txt file to save the dump data
	err := file.CreateFile(filepath.Join(dumpWorkingPath, taskId+".txt"))
	if err != nil {
		errHandler(false, "create file error: "+err.Error())
		return task.FAIL, err
	}

	// save the root first
	firstQuery := `{
		"query": {
			"bool": {
				"must": [
					{ "term": { "agent": "` + key + `" }},
					{ "match_phrase": { "etc_main": "` + hitPath + `" }},
					{ "term": { "explorer.fileId": ` + fileId + ` }}
				]
			}
		}
	}`

	hitsArray := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", firstQuery, "uuid", 0)
	if len(hitsArray) == 0 {
		errHandler(true, "hitsArray is empty")
		return task.FAIL, errors.New("hitsArray is empty")
	}

	hitMap, ok := hitsArray[0].(map[string]interface{})
	if !ok {
		errHandler(true, "hit is not a map")
		return task.FAIL, errors.New("hit is not a map")
	}

	source, ok := hitMap["_source"].(map[string]interface{})
	if !ok {
		errHandler(true, "hitMap[\"_source\"] is not a map")
		return task.FAIL, errors.New("hitMap[\"_source\"] is not a map")
	}

	explorerData, ok := source["explorer"].(map[string]interface{})
	if !ok {
		errHandler(true, "source[\"explorer\"] is not a map")
		return task.FAIL, errors.New("source[\"explorer\"] is not a map")
	}

	// save in the txt file
	var data string
	if strings.Split(explorerData["disk"].(string), "|")[1] == "NTFS" {
		// path|isDirectory|fileId
		if explorerData["isDirectory"].(bool) {
			data = explorerData["path"].(string) + "|1|" +
				strconv.FormatFloat(explorerData["fileId"].(float64), 'f', 0, 64) + "\n"
		} else {
			data = explorerData["path"].(string) + "|0|" +
				strconv.FormatFloat(explorerData["fileId"].(float64), 'f', 0, 64) + "\n"
		}
	} else if strings.Split(explorerData["disk"].(string), "|")[1] == "FAT32" {
		// path|isDirectory|startCluster|fileSize
		if explorerData["isDirectory"].(bool) {
			data = explorerData["path"].(string) + "|1|" +
				strconv.FormatFloat(explorerData["startCluster"].(float64), 'f', 0, 64) + "|" +
				strconv.FormatFloat(explorerData["dataLen"].(float64), 'f', 0, 64) + "\n"
		} else {
			data = explorerData["path"].(string) + "|0|" +
				strconv.FormatFloat(explorerData["startCluster"].(float64), 'f', 0, 64) + "|" +
				strconv.FormatFloat(explorerData["dataLen"].(float64), 'f', 0, 64) + "\n"
		}
	} else {
		// path|isDirectory
		if explorerData["isDirectory"].(bool) {
			data = explorerData["path"].(string) + "|1\n"
		} else {
			data = explorerData["path"].(string) + "|0\n"
		}
	}
	file.WriteFile(filepath.Join(dumpWorkingPath, taskId+".txt"), []byte(data))

	// retrive all the file info that need to dump from elastic and save in a txt file
	query := `{
		"query": {
		  	"bool": {
				"must": [
			  		{
						"query_string": {
				  			"query": "` + followingPath + `*"
						}
			  		}
				],
				"filter": [
			  		{ "term": { "agent": "` + key + `" }}
				]
		  	}
		}
	}`

	hitsArray = elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", query, "uuid", 0)
	if len(hitsArray) == 0 {
		errHandler(true, "hitsArray is empty")
		return task.FAIL, errors.New("hitsArray is empty")
	}

	isFirst, isDirectory := true, true
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			errHandler(true, "hit is not a map")
			return task.FAIL, errors.New("hit is not a map")
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			errHandler(true, "hitMap[\"_source\"] is not a map")
			return task.FAIL, errors.New("hitMap[\"_source\"] is not a map")
		}

		// skip the wrong path and original path
		if !strings.Contains(source["etc_main"].(string), filePath) || source["etc_main"].(string) == filePath {
			continue
		}

		explorerData, ok := source["explorer"].(map[string]interface{})
		if !ok {
			errHandler(true, "source[\"explorer\"] is not a map")
			return task.FAIL, errors.New("source[\"explorer\"] is not a map")
		}

		// save in the txt file
		var data string
		if strings.Split(explorerData["disk"].(string), "|")[1] == "NTFS" {
			// path|isDirectory|fileId
			if explorerData["isDirectory"].(bool) {
				data = explorerData["path"].(string) + "|1|" +
					strconv.FormatFloat(explorerData["fileId"].(float64), 'f', 0, 64) + "\n"
			} else {
				data = explorerData["path"].(string) + "|0|" +
					strconv.FormatFloat(explorerData["fileId"].(float64), 'f', 0, 64) + "\n"
			}
		} else if strings.Split(explorerData["disk"].(string), "|")[1] == "FAT32" {
			// path|isDirectory|startCluster|fileSize
			if explorerData["isDirectory"].(bool) {
				data = explorerData["path"].(string) + "|1|" +
					strconv.FormatFloat(explorerData["startCluster"].(float64), 'f', 0, 64) + "|" +
					strconv.FormatFloat(explorerData["dataLen"].(float64), 'f', 0, 64) + "\n"
			} else {
				data = explorerData["path"].(string) + "|0|" +
					strconv.FormatFloat(explorerData["startCluster"].(float64), 'f', 0, 64) + "|" +
					strconv.FormatFloat(explorerData["dataLen"].(float64), 'f', 0, 64) + "\n"
			}
		} else {
			// path|isDirectory
			if explorerData["isDirectory"].(bool) {
				data = explorerData["path"].(string) + "|1\n"
			} else {
				data = explorerData["path"].(string) + "|0\n"
			}
		}
		file.WriteFile(filepath.Join(dumpWorkingPath, taskId+".txt"), []byte(data))

		if isFirst {
			isDirectory = explorerData["isDirectory"].(bool)
			isFirst = false
		}
	}

	toAgentMsg := filePath
	if isDirectory {
		toAgentMsg += "|1"
	} else {
		toAgentMsg += "|0"
	}

	err = clientsearchsend.SendUserTCPtoClientUsingKey(key, task.GET_DUMP_DRIVE, toAgentMsg)
	if err != nil {
		errHandler(true, "send tcp error: "+err.Error())
		return task.FAIL, err
	}

	// update progress
	request.LoadDumpReady(request.ReadyData{
		TaskId:   taskId,
		Progress: config.Viper.GetInt("DUMP_DRIVE_FIRST_PART"),
	})
	logger.Info(key + "::StartDumpDrive: update progress to " + strconv.Itoa(config.Viper.GetInt("DUMP_DRIVE_FIRST_PART")))

	return task.SUCCESS, nil
}

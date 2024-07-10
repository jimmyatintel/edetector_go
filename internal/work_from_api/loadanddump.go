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
	"strings"
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
			Progress: -1,
		})
		return task.FAIL, err
	}

	err := clientsearchsend.SendUserTCPtoClientUsingKey(key, task.GET_LOAD_DLL, msg)
	if err != nil {
		logger.Error("StartLoadDll: send tcp error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
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
			Progress: -1,
		})
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func StartDumpDrive(key, taskId, msg string) (task.TaskResult, error) {
	logger.Info(key + "::StartDumpDrive: " + msg)
	msgs := strings.Split(msg, "|")
	_, filePath := msgs[0], msgs[1]

	// open a txt file to save the dump data
	err := file.CreateFile(filepath.Join(dumpWorkingPath, taskId+".txt"))
	if err != nil {
		logger.Error("StartDumpDrive: create file error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			Progress: -1,
		})
		return task.FAIL, err
	}

	// retrive all the file info that need to dump from elastic and save in a txt file
	query := `{
			"query": {
				"match": {
					"agent": "` + key + `",
					"explorer.path": ` + filePath + `
				}
			}
		}`
	hitsArray := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", query, "uuid", 0)
	hitMap, ok := hitsArray[0].(map[string]interface{})
	if !ok {
		os.Remove(filepath.Join(dumpWorkingPath, taskId+".txt"))
		logger.Error("StartDumpDrive: hit is not a map")
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			Progress: -1,
		})
		return task.FAIL, errors.New("hit is not a map")
	}

	source, ok := hitMap["_source"].(map[string]interface{})
	if !ok {
		os.Remove(filepath.Join(dumpWorkingPath, taskId+".txt"))
		logger.Error("StartDumpDrive: hitMap[\"_source\"] is not a map")
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			Progress: -1,
		})
		return task.FAIL, errors.New("hitMap[\"_source\"] is not a map")
	}

	explorerData, ok := source["explorer"].(map[string]interface{})
	if !ok {
		os.Remove(filepath.Join(dumpWorkingPath, taskId+".txt"))
		logger.Error("StartDumpDrive: source[\"explorer\"] is not a map")
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			Progress: -1,
		})
		return task.FAIL, errors.New("source[\"explorer\"] is not a map")
	}

	// check if the data is too large, smaller than 2GB is ok
	if explorerData["dataLen"].(int) > 2*1024*1024*1024 {
		os.Remove(filepath.Join(dumpWorkingPath, taskId+".txt"))
		logger.Error("StartDumpDrive: data is larger than 2GB")
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			Progress: -1,
		})
		return task.FAIL, errors.New("data is larger than 2GB")
	}

	uuid := source["uuid"].(string)
	findAndSaveDumpPaths(taskId, []string{uuid})

	// save task id in redis
	redisKey := key + string(task.START_DUMP_DRIVE) + filePath
	err = redis.RedisSet(redisKey, taskId)
	if err != nil {
		os.Remove(filepath.Join(dumpWorkingPath, taskId+".txt"))
		logger.Error("redis set error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			Progress: -1,
		})
		return task.FAIL, err
	}

	err = clientsearchsend.SendUserTCPtoClientUsingKey(key, task.GET_DUMP_DRIVE, strings.Split(msg, "|")[1])
	if err != nil {
		os.Remove(filepath.Join(dumpWorkingPath, taskId+".txt"))
		logger.Error("StartDumpDrive: send tcp error" + err.Error())
		request.LoadDumpReady(request.ReadyData{
			TaskId:   taskId,
			Failed:   "InternalServerError",
			Progress: -1,
		})
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func findAndSaveDumpPaths(taskId string, uuids []string) {
	for _, uuid := range uuids {
		query := `{
			"query": {
				"match": {
					"uuid": "` + uuid + `"
				}
			}
		}`
		hitsArray := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", query, "uuid", 0)

		for _, hit := range hitsArray {
			hitMap, ok := hit.(map[string]interface{})
			if !ok {
				logger.Error("FindAndSaveDumpPaths: hit is not a map")
				continue
			}
			source, ok := hitMap["_source"].(map[string]interface{})
			if !ok {
				logger.Error("FindAndSaveDumpPaths: hitMap[\"_source\"] is not a map")
				continue
			}
			explorerData, ok := source["explorer"].(map[string]interface{})
			if !ok {
				logger.Error("FindAndSaveDumpPaths: source[\"explorer\"] is not a map")
				continue
			}

			// save in the txt file
			var data string
			if strings.Split(explorerData["disk"].(string), "|")[1] == "NTFS" {
				// path|fileId
				data = explorerData["path"].(string) + "|" + explorerData["fileId"].(string) + "\n"
			} else if strings.Split(explorerData["disk"].(string), "|")[1] == "FAT32" {
				// path|startCluster|fileSize
				data = explorerData["path"].(string) + "|" + explorerData["startCluster"].(string) + "|" + explorerData["fileSize"].(string) + "\n"
			} else {
				// path
				data = explorerData["path"].(string) + "\n"
			}
			file.WriteFile(filepath.Join(dumpWorkingPath, taskId+".txt"), []byte(data))

			// iterate the children
			children, ok := source["children"].([]string)
			if !ok {
				logger.Error("FindAndSaveDumpPaths: source[\"children\"] is not an array")
				continue
			}

			if len(children) > 0 {
				findAndSaveDumpPaths(taskId, children)
			}
		}
	}
}

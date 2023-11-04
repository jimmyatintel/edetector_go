package main

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
)

// import (
// 	"edetector_go/config"
// 	"edetector_go/pkg/elastic"
// 	"edetector_go/pkg/logger"
// )

func main() {
	vp, err := config.LoadConfig()
	if vp == nil {
		panic(err)
	}
	if true {
		logger.InitLogger("cmd/test.txt", "test", "TEST")
		logger.Info("Logger is enabled please check all out info in log file")
	}
	connString, err := mariadb.Connect_init()
	if err != nil {
		logger.Panic("Error connecting to mariadb: " + err.Error())
		panic(err)
	} else {
		logger.Info("Mariadb connectionString: " + connString)
	}
	if true {
		if db := redis.Redis_init(); db == nil {
			logger.Panic("Error connecting to redis")
			panic(err)
		}
	}
	if true {
		rabbitmq.Rabbit_init()
		logger.Info("Rabbit is enabled.")
	}
	if true {
		elastic.Elastic_init()
		logger.Info("Elastic is enabled.")
	}

	// unstagePath := filepath.Join("scanUnstage", ("0a741ef4fe8747178b7076e1e0161e04.txt"))
	// err = work.ParseScan(unstagePath, "0a741ef4fe8747178b7076e1e0161e04")
	// if err != nil {
	// 	logger.Error("Failed to parse" + err.Error())
	// }

	// 	index := config.Viper.GetString("ELASTIC_PREFIX") + "_explorer_relation"

	//	searchQuery := `{
	//		"query": {
	//			"bool": {
	//			  "must": [
	//				{ "query_string": { "fields": ["agent"], "query": "556c6050a5204cc8b0dc39bf9d15e941" } },
	//				{ "query_string": { "fields": ["parent"], "query": "b6b75a86-8178-4827-b15e-0ed75f57eb42" } }
	//			  ]
	//			}
	//		  }
	//		}`
	//
	// hitsArray := elastic.SearchRequest(index, searchQuery)
	//
	//	for _, hit := range hitsArray {
	//		hitMap, ok := hit.(map[string]interface{})
	//		if !ok {
	//			logger.Error("Hit is not a map")
	//			continue
	//		}
	//		docID, ok := hitMap["_id"].(string)
	//		if !ok {
	//			logger.Error("docID not found")
	//			continue
	//		}
	//		_, err = elastic.UpdateByDocIDRequest(index, docID, "HIHI")
	//		if err != nil {
	//			logger.Error("Error updating child: " + err.Error())
	//		}
	//	}
}

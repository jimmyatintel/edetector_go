package main

// import (
// 	"edetector_go/config"
// 	"edetector_go/pkg/elastic"
// 	"edetector_go/pkg/logger"
// )

func main() {
	// 	vp, err := config.LoadConfig()
	// 	if vp == nil {
	// 		panic(err)
	// 	}
	// 	if true {
	// 		logger.InitLogger("cmd/test.txt", "test", "TEST")
	// 		logger.Info("Logger is enabled please check all out info in log file")
	// 	}
	// 	elastic.Elastic_init()

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

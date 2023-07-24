package elasticquery

import (
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

func Send_to_main_elastic(index string, template mainSource) {
	request, err := template.Elastical()
	if err != nil {
		logger.Error("Error sending to elastic: ", zap.Any("error", err.Error()))
		return
	}
	fmt.Println("Sending to Elastic:", string(index), string(request))
	elastic.IndexRequest(index, string(request))
}

// func Send_to_sub_elastic(index string, template mainSource) {
// 	request, err := template.Elastical()
// 	if err != nil {
// 		logger.Error("Error sending to elastic: ", zap.Any("error", err.Error()))
// 		return
// 	}
// 	fmt.Println("Sending to Elastic:", string(index), string(request))
// 	elastic.IndexRequest(index, string(request))
// }
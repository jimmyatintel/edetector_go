package elasticquery

import (
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
)

func Send_to_elastic(template source, data []Request_data) {
		for _, v := range data {
			template.Data = v
			request, err := template.Elastical()
			if err != nil {
				logger.Error(err.Error())
			}
			elastic.IndexRequest("network", string(request))
		}
}

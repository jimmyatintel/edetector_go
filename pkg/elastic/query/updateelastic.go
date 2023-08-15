package elasticquery

import (
	"edetector_go/config"
	elastic "edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"strings"

	"go.uber.org/zap"
)

func UpdateNetworkInfo(agent string, networkSet map[string]struct{}) {
	for key := range networkSet {
		values := strings.Split(key, ",")
		id := values[0]
		time := values[1]
		err := elastic.UpdateRequest(agent, id, time, config.Viper.GetString("ELASTIC_PREFIX")+"_memory")
		if err != nil {
			logger.Error("Error updating detect process: ", zap.Any("message", err.Error()))
		}
	}
}

package elasticquery

import (
	"edetector_go/config"
	elastic "edetector_go/pkg/elastic"
	"strings"
)

func UpdateNetworkInfo(agent string, networkSet map[string]struct{}) {
	for key := range networkSet {
		values := strings.Split(key, ",")
		id := values[0]
		time := values[1]
		elastic.UpdateRequest(agent, id, time, config.Viper.GetString("ELASTIC_PREFIX")+"_memory")

	}
}

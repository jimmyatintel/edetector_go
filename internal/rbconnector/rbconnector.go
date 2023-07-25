package rbconnector

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type message struct {
	Index string `json:"index"`
	Data  string `json:"data"`
}

func init() {
	// rbconnector.init()
	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		fmt.Println("Error loading feature flag")
		return
	}
	vp := config.LoadConfig()
	if vp == nil {
		fmt.Println("Error loading config file")
		return
	}
}

func Start() {
	// rbconnector.Start()
	rabbitmq.Rabbit_init()
	rabbitmq.Declare("ed_low")
	rabbitmq.Declare("ed_mid")
	rabbitmq.Declare("ed_high")
	go low_speed()
	go mid_speed()
	go high_speed()
}
func high_speed() {
	msgs, err := rabbitmq.Consume("ed_high")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CONNECT TO high SPEED QUEUE")
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
		var m message
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		err = elastic.IndexRequest(m.Index, m.Data)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
}

func mid_speed() {
	msgs, err := rabbitmq.Consume("ed_mid")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CONNECT TO mid SPEED QUEUE")
	last_send := time.Now()
	var bulkdata []string
	var bulkaction []string
	for msg := range msgs {
		var m message
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		bulkdata = append(bulkdata, m.Data)
		bulkaction = append(bulkaction, fmt.Sprintf(`{ "index" : { "_index" : "%s", "_type" : "_doc" } }`, m.Index))
		if time.Since(last_send) > time.Duration(config.Viper.GetInt("MID_TUNNEL_TIME"))*time.Second || len(bulkaction) > config.Viper.GetInt("MID_TUNNEL_SIZE") {
			last_send = time.Now()
			err = elastic.BulkIndexRequest(m.Index, bulkaction, bulkdata)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			bulkdata = nil
			bulkaction = nil
		}
	}

}
func low_speed() {
	msgs, err := rabbitmq.Consume("ed_low")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CONNECT TO LOW SPEED QUEUE")
	last_send := time.Now()
	var bulkdata []string
	var bulkaction []string
	for msg := range msgs {
		var m message
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		bulkdata = append(bulkdata, m.Data)
		bulkaction = append(bulkaction, fmt.Sprintf(`{ "index" : { "_index" : "%s", "_type" : "_doc" } }`, m.Index))
		if time.Since(last_send) > time.Duration(config.Viper.GetInt("LOW_TUNNEL_TIME"))*time.Second || len(bulkaction) > config.Viper.GetInt("LOW_TUNNEL_SIZE") {
			last_send = time.Now()
			err = elastic.BulkIndexRequest(m.Index, bulkaction, bulkdata)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			bulkdata = nil
			bulkaction = nil
		}
	}
}

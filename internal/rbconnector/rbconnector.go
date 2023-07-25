package rbconnector

import (
	"context"
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Message struct {
	Index string `json:"index"`
	Data  string `json:"data"`
}

func init() {
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
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("CONNECTOR_LOG_FILE"))
		fmt.Println("logger is enabled please check all out info in log file: ", config.Viper.GetString("CONNECTOR_LOG_FILE"))
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
		err := elastic.SetElkClient()
		if err != nil {
			logger.Error("Error connecting to elastic: " + err.Error())
		}
		fmt.Println("elastic is enabled.")
	}
}

func Start() {
	// rbconnector.Start()
	Quit := make(chan os.Signal, 1)
	_, cancel := context.WithCancel(context.Background())
	rabbitmq.Rabbit_init()
	rabbitmq.Declare("ed_low")
	rabbitmq.Declare("ed_mid")
	rabbitmq.Declare("ed_high")
	go low_speed()
	go mid_speed()
	go high_speed()
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	cancel()
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
		var m Message
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
		var m Message
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
		var m Message
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

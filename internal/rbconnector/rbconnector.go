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
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Message struct {
	Index string `json:"index"`
	Data  string `json:"data"`
}

var mid_mutex *sync.Mutex
var low_mutex *sync.Mutex

var mid_bulkdata []string
var mid_bulkaction []string
var low_bulkdata []string
var low_bulkaction []string

func connector_init() {
	mid_mutex = &sync.Mutex{}
	low_mutex = &sync.Mutex{}
	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Error("Error loading feature flag")
		return
	}
	vp := config.LoadConfig()
	if vp == nil {
		logger.Error("Error loading config file")
		return
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("CONNECTOR_LOG_FILE"))
		logger.Info("logger is enabled please check all out info in log file: ", zap.Any("message", config.Viper.GetString("CONNECTOR_LOG_FILE")))
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
		err := elastic.SetElkClient()
		if err != nil {
			logger.Error("Error connecting to elastic: " + err.Error())
		}
		logger.Info("elastic is enabled.")
	}
	fmt.Println("hihi213")
}

func Start() {
	connector_init()
	Quit := make(chan os.Signal, 1)
	_, cancel := context.WithCancel(context.Background())
	rabbitmq.Rabbit_init()
	rabbitmq.Declare("ed_low")
	rabbitmq.Declare("ed_mid")
	rabbitmq.Declare("ed_high")
	fmt.Println("123")
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
		logger.Error("High speed consumer error: " + err.Error())
		return
	}
	logger.Info("Connected to high speed queue")
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
		logger.Error("Mid speed consumer error: " + err.Error())
		return
	}
	logger.Info("Connected to mid speed queue")
	go count_timer(config.Viper.GetInt("MID_TUNNEL_TIME"), config.Viper.GetInt("MID_TUNNEL_SIZE"), &mid_bulkaction, &mid_bulkdata, mid_mutex)
	for msg := range msgs {
		var m Message
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		mid_mutex.Lock()
		mid_bulkdata = append(mid_bulkdata, m.Data)
		mid_bulkaction = append(mid_bulkaction, fmt.Sprintf(`{ "index" : { "_index" : "%s", "_type" : "_doc" } }`, m.Index))
		mid_mutex.Unlock()
		for len(mid_bulkaction) > config.Viper.GetInt("MID_TUNNEL_SIZE") {
			time.Sleep(1 * time.Second)
		}
	}
}

func low_speed() {
	msgs, err := rabbitmq.Consume("ed_low")
	if err != nil {
		logger.Error("Low speed consumer error: " + err.Error())
		return
	}
	logger.Info("Connected to low speed queue")
	go count_timer(config.Viper.GetInt("LOW_TUNNEL_TIME"), config.Viper.GetInt("LOW_TUNNEL_SIZE"), &low_bulkaction, &low_bulkdata, low_mutex)
	for msg := range msgs {
		var m Message
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		low_mutex.Lock()
		low_bulkdata = append(low_bulkdata, m.Data)
		low_bulkaction = append(low_bulkaction, fmt.Sprintf(`{ "index" : { "_index" : "%s", "_type" : "_doc" } }`, m.Index))
		low_mutex.Unlock()
		for len(low_bulkaction) > config.Viper.GetInt("LOW_TUNNEL_SIZE") {
			time.Sleep(1 * time.Second)
		}
	}
}

func count_timer(tunnel_time int, size int, bulkaction *[]string, bulkdata *[]string, mutex *sync.Mutex) {
	logger.Info("Counting timer started")
	last_send := time.Now()
	for {
		mutex.Lock()
		if ((time.Since(last_send) > time.Duration(tunnel_time)*time.Second) && len(*bulkaction) > 0) || len(*bulkaction) > size {
			err := elastic.BulkIndexRequest(*bulkaction, *bulkdata)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			*bulkdata = nil
			*bulkaction = nil
			last_send = time.Now()
		}
		mutex.Unlock()
		time.Sleep(3 * time.Second)
	}
}

package rbconnector

//TBD
import (
	"context"
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	elaInsert "edetector_go/pkg/elastic/insert"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

var mid_mutex *sync.Mutex
var low_mutex *sync.Mutex

var mid_bulkdata []string
var mid_bulkaction []string

var consumer_count map[string]int = make(map[string]int)
var queues = []string{"ed_low_collect", "ed_low_explorer"}

var hostname string
var port string
var username string
var password string
var vhost string

type RateDetails struct {
	Rate float64 `json:"rate"`
}

type MessageStats struct {
	PublishDetails    RateDetails `json:"publish_details"`
	DeliverGetDetails RateDetails `json:"deliver_details"`
}

type QueueInfo struct {
	MessageStats MessageStats `json:"message_stats"`
}

func init() {
	mid_mutex = &sync.Mutex{}
	low_mutex = &sync.Mutex{}

	// fflag.Get_fflag()
	// if fflag.FFLAG == nil {
	// 	logger.Panic("Error loading feature flag")
	// 	panic("error loading feature flag")
	// }
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Panic("Error loading config file: " + err.Error())
		panic(err)
	}
	if true {
		logger.InitLogger(config.Viper.GetString("CONNECTOR_LOG_FILE"), "connector", "CONNECTOR")
		logger.Info("Logger is enabled please check all out info in log file: " + config.Viper.GetString("CONNECTOR_LOG_FILE"))
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
		elastic.Elastic_init()
		logger.Info("Elastic is enabled.")
	}
	hostname = config.Viper.GetString("RABBITMQ_IP")
	port = config.Viper.GetString("RABBITMQ_WEB_PORT")
	username = config.Viper.GetString("RABBITMQ_USERNAME")
	password = config.Viper.GetString("RABBITMQ_PASSWORD")
	vhost = url.PathEscape(config.Viper.GetString("RABBITMQ_VHOST"))
	consumer_count["ed_low_collect"] = 0
	consumer_count["ed_low_explorer"] = 0
}

func Start(version string) {
	logger.Info("Welcome to edetector connector: " + version)
	Quit := make(chan os.Signal, 1)
	_, cancel := context.WithCancel(context.Background())
	rabbitmq.Rabbit_init()
	rabbitmq.Declare("ed_low_collect")
	rabbitmq.Declare("ed_low_explorer")
	rabbitmq.Declare("ed_mid")
	rabbitmq.Declare("ed_high")
	for _, queue := range queues {
		consumer_count[queue]++
		time.Sleep(1 * time.Second)
		go low_speed(queue, consumer_count[queue])
	}
	time.Sleep(1 * time.Second)
	go mid_speed()
	time.Sleep(1 * time.Second)
	go high_speed()
	go scale()
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	cancel()
}

func scale() {
	logger.Info("Scale check started")
	for {
		time.Sleep(time.Duration(config.Viper.GetInt("RABBITMQ_SCALE_SLEEP")) * time.Second)
		for _, queue := range queues {
			count, err := GetMessageCount(queue)
			if err != nil {
				logger.Error("Error getting rabbit messages: " + err.Error())
				continue
			}
			if (count > config.Viper.GetInt("RABBITMQ_SCALE_THRESHOLD") || PublishFasterThanDeliver(queue)) && consumer_count[queue] < config.Viper.GetInt("RABBITMQ_MAX_CONSUMER") {
				consumer_count[queue]++
				time.Sleep(1 * time.Second)
				go low_speed(queue, consumer_count[queue]) // add a consumer
			} else if consumer_count[queue] != 1 {
				rabbitmq.Cancel(queue + "-" + strconv.Itoa(consumer_count[queue])) // cancel the consumer
				if err != nil {
					logger.Error("Failed to cancel consumer" + err.Error())
				}
				logger.Info("Removed a consumer: " + queue + "-" + strconv.Itoa(consumer_count[queue]))
				consumer_count[queue]--
			}
		}
	}
}

func high_speed() {
	var msgs <-chan amqp.Delivery
	var err error
	for {
		msgs, err = rabbitmq.Consume("ed_high", 1, config.Viper.GetInt("LOW_TUNNEL_SIZE"))
		if err != nil {
			logger.Error("High speed consumer error: " + err.Error())
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	logger.Info("Added a consumer: ed_high")
	for msg := range msgs {
		logger.Info("Received a message: " + string(msg.Body))
		var m rabbitmq.Message
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			logger.Error("Error unmarshaling: " + err.Error())
			continue
		}
		err = elastic.IndexRequest(m.Index, m.Data, 0)
		if err != nil {
			logger.Error("Index request error: " + err.Error())
			continue
		}
		msg.Ack(false)
	}
}

func mid_speed() {
	var msgs <-chan amqp.Delivery
	var err error
	for {
		msgs, err = rabbitmq.Consume("ed_mid", 1, config.Viper.GetInt("LOW_TUNNEL_SIZE"))
		if err != nil {
			logger.Error("Mid speed consumer error: " + err.Error())
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	logger.Info("Added a consumer: ed_mid")
	go count_timer(config.Viper.GetInt("MID_TUNNEL_TIME"), config.Viper.GetInt("MID_TUNNEL_SIZE"), &mid_bulkaction, &mid_bulkdata, mid_mutex)
	for msg := range msgs {
		var m rabbitmq.Message
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			logger.Error("Error unmarshaling: " + err.Error())
			continue
		}
		mid_mutex.Lock()
		mid_bulkdata = append(mid_bulkdata, m.Data)
		mid_bulkaction = append(mid_bulkaction, fmt.Sprintf(`{ "index" : { "_index" : "%s" } }`, m.Index))
		mid_mutex.Unlock()
		msg.Ack(false)
		for len(mid_bulkaction) > (config.Viper.GetInt("MID_TUNNEL_SIZE") * 2) {
			time.Sleep(2 * time.Second)
		}
	}
}

func low_speed(queue string, count int) {
	var msgs <-chan amqp.Delivery
	var err error
	for {
		msgs, err = rabbitmq.Consume(queue, count, config.Viper.GetInt("LOW_TUNNEL_SIZE"))
		if err != nil {
			logger.Error("Low speed consumer error: " + err.Error())
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	logger.Info("Added a consumer: " + queue + "-" + strconv.Itoa(count))
	var low_bulkdata []string
	var low_bulkaction []string
	go count_timer(config.Viper.GetInt("LOW_TUNNEL_TIME"), config.Viper.GetInt("LOW_TUNNEL_SIZE"), &low_bulkaction, &low_bulkdata, low_mutex)
	for msg := range msgs {
		var m rabbitmq.Message
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			logger.Error("Error unmarshaling: " + err.Error())
			continue
		}
		low_mutex.Lock()
		low_bulkdata = append(low_bulkdata, m.Data)
		low_bulkaction = append(low_bulkaction, fmt.Sprintf(`{ "index" : { "_index" : "%s" } }`, m.Index))
		low_mutex.Unlock()
		msg.Ack(false)
		for len(low_bulkaction) > (config.Viper.GetInt("LOW_TUNNEL_SIZE") * 2) {
			time.Sleep(2 * time.Second)
		}
	}
}

func count_timer(tunnel_time int, size int, bulkaction *[]string, bulkdata *[]string, mutex *sync.Mutex) {
	logger.Info("Counting timer started")
	last_send := time.Now()
	for {
		mutex.Lock()
		if ((time.Since(last_send) > time.Duration(tunnel_time)*time.Second) && len(*bulkaction) > 0) || len(*bulkaction) > size {
			err := elaInsert.BulkInsert(*bulkaction, *bulkdata)
			if err != nil {
				logger.Error("Bulk index request error: " + err.Error())
				time.Sleep(10 * time.Second)
			} else {
				*bulkdata = nil
				*bulkaction = nil
				last_send = time.Now()
			}
		}
		mutex.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func GetMessageCount(queue string) (int, error) {
	path := fmt.Sprintf("http://%s:%s/api/queues/%s/%s", hostname, port, vhost, queue)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, errors.New("error getting rabbitmq queue info")
	}
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return 0, err
	}
	count, ok := result["messages"].(float64)
	if !ok {
		return 0, fmt.Errorf("messages not found in response")
	}
	return int(count), nil
}

func PublishFasterThanDeliver(queue string) bool {
	path := fmt.Sprintf("http://%s:%s/api/queues/%s/%s", hostname, port, vhost, queue)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		logger.Error("Error creating new request:" + err.Error())
		return false
	}
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error getting rabbitmq queue info:" + err.Error())
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		logger.Error("Error getting rabbitmq queue info: " + response.Status)
		return false
	}
	var queueInfo QueueInfo
	if err := json.NewDecoder(response.Body).Decode(&queueInfo); err != nil {
		logger.Error("Error decoding rabbitmq queue info:" + err.Error())
		return false
	}
	logger.Debug("Queue: " + queue +
		" ,PublishRate: " + strconv.FormatFloat(queueInfo.MessageStats.PublishDetails.Rate, 'f', -1, 64) +
		" ,DeliverGetRate: " + strconv.FormatFloat(queueInfo.MessageStats.DeliverGetDetails.Rate, 'f', -1, 64))
	return queueInfo.MessageStats.PublishDetails.Rate > queueInfo.MessageStats.DeliverGetDetails.Rate
}

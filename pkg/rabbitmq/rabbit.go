package rabbitmq

import (
	"edetector_go/config"
	"edetector_go/pkg/logger"
	"strconv"

	"errors"

	"github.com/streadway/amqp"
)

func NewRabbitMQ(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}

func NewQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
}

var Connection *amqp.Connection
var channel *amqp.Channel

func Rabbit_init() {
	// ...
	var err error
	hostname := config.Viper.GetString("RABBITMQ_IP")
	port := config.Viper.GetString("RABBITMQ_PORT")
	username := config.Viper.GetString("RABBITMQ_USERNAME")
	password := config.Viper.GetString("RABBITMQ_PASSWORD")
	Url := "amqp://" + username + ":" + password + "@" + hostname + ":" + port + "/"
	Connection, err = NewRabbitMQ(Url)
	if err != nil {
		logger.Panic("Failed to connect to RabbitMQ")
		panic(err)
	}
	channel, err = Connection.Channel()
	if err != nil {
		logger.Panic("Failed to connect to RabbitMQ")
		panic(err)
	}
	go func() {
		notifyClose := channel.NotifyClose(make(chan *amqp.Error))
		err := <-notifyClose
		if err != nil {
			logger.Error("Connection closed: " + err.Error())
			Rabbit_init()
		}
	}()
}

func Declare(name string) (amqp.Queue, error) {
	if channel == nil {
		return amqp.Queue{}, errors.New("failed to declare queue: channel is nil")
	}
	return NewQueue(channel, name)
}

func Publish(queue string, body []byte) error {
	if channel == nil {
		return errors.New("failed to publish message: channel is nil")
	}
	return channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func Consume(queue string, id int, count int) (<-chan amqp.Delivery, error) {
	if channel == nil {
		return nil, errors.New("failed to consume message: channel is nil")
	}
	err := channel.Qos(count, 0, false)
	if err != nil {
		logger.Error("Error setting consume messages")
	}
	name := queue + "-" + strconv.Itoa(id)
	return channel.Consume(
		queue,
		name,
		false,
		false,
		false,
		false,
		nil,
	)
}

func Cancel(consumer string) {
	if channel == nil {
		return
	}
	channel.Cancel(consumer, false)
}

func Connection_close() {
	if Connection != nil {
		Connection.Close()
	}
}

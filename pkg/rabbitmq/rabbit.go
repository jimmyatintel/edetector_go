package rabbitmq

import (
	"edetector_go/pkg/logger"
	"os"

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
	hostname := os.Getenv("RABBITMQ_IP")
	port := os.Getenv("RABBITMQ_PORT")
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	Url := "amqp://" + username + ":" + password + "@" + hostname + ":" + port + "/"
	Connection, err = NewRabbitMQ(Url)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ")
	}
	channel, err = Connection.Channel()
	if err != nil {
		panic(err)
	}
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
func Consume(queue string) (<-chan amqp.Delivery, error) {
	if channel == nil {
		return nil, errors.New("failed to consume message: channel is nil")
	}
	return channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
func Connection_close() {
	if Connection != nil {
		Connection.Close()
	}
}

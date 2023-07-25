package rbconnector

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/rabbitmq"
	"fmt"
	"log"
)

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
	}
}
func mid_speed() {
	msgs, err := rabbitmq.Consume("ed_mid")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CONNECT TO mid SPEED QUEUE")
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
	}

}
func low_speed() {
	msgs, err := rabbitmq.Consume("ed_low")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CONNECT TO LOW SPEED QUEUE")
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
	}
}

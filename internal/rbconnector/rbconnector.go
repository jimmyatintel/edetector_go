package rbconnector

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/rabbitmq"
	"fmt"
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

}

package main

import (
	"edetector_go/config"
	"edetector_go/pkg/redis"
	"fmt"
)

func main() {
	vp, err := config.LoadConfig()
	if vp == nil {
		panic(err)
	}
	if db := redis.Redis_init(); db == nil {
		panic(err)
	}
	length := 32
	keys := redis.GetKeysByLength(length)
	values := redis.GetValuesForKeys(keys)

	for key, value := range values {
		fmt.Println("Key:", key, "Value:", value)
	}
}

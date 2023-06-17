package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

// Config is a struct that holds the configuration for the application
func LoadConfig() *viper.Viper {
	vp := viper.New()
	vp.SetConfigName("app")
	vp.SetConfigType("env")
	vp.AddConfigPath("../config")
	vp.AutomaticEnv()
	if err := vp.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", vp.ConfigFileUsed())
		Viper = vp
		return vp
	} else {
		fmt.Println("Error loading config file:", err)
		return nil
	}
}

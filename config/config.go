package config

import (
	"github.com/spf13/viper"
)

var Viper *viper.Viper

// Config is a struct that holds the configuration for the application
func LoadConfig() (*viper.Viper, error) {
	vp := viper.New()
	vp.SetConfigName("app")
	vp.SetConfigType("env")
	vp.AddConfigPath("config")
	vp.AutomaticEnv()
	// if err := vp.ReadInConfig(); err != nil {
	// 	return nil, err
	// }
	Viper = vp
	return vp, nil
}
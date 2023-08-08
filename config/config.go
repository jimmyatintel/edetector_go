package config

import (
	"edetector_go/pkg/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Viper *viper.Viper

// Config is a struct that holds the configuration for the application
func LoadConfig() *viper.Viper {
	vp := viper.New()
	// vp.SetConfigName("app")
	// vp.SetConfigType("env")
	// vp.AddConfigPath("config")
	vp.AutomaticEnv()
	if err := vp.ReadInConfig(); err == nil {
		logger.Debug("Using config file:", zap.Any("config", vp.ConfigFileUsed()))
		Viper = vp
		return vp
	} else {
		logger.Error("Error loading config file:", zap.Any("error", err))
		return nil
	}
}

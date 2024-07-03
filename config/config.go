package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

// Config is a struct that holds the configuration for the application
func LoadConfig() (*viper.Viper, error) {
	vp := viper.New()
	vp.AutomaticEnv() // Automatically read from environment variables

	// Define and parse the command-line flag for environment
	env := flag.String("env", "prod", "Environment: dev, prod")
	flag.Parse()

	// Validate environment value
	if *env != "dev" && *env != "prod" {
		return nil, fmt.Errorf("invalid environment: %s", *env)
	}

	// env = dev -> read from config file
	if *env == "dev" {
		vp.SetConfigName("app")
		vp.SetConfigType("env")
		vp.AddConfigPath("config")
		if err := vp.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	Viper = vp
	return vp, nil
}

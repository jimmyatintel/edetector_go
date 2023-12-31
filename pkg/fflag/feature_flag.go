package fflag

import (
	"edetector_go/pkg/logger"

	flagsmith "github.com/Flagsmith/flagsmith-go-client"
	"github.com/Unleash/unleash-client-go/v3"
)

type metricsInterface struct {
}

// Initialise the Flagsmith client
var FFLAG *flagsmith.Client

func init_from_gitlab() {
	err := unleash.Initialize(
		unleash.WithUrl("https://git.chainsecurity.local/api/v4/feature_flags/unleash/63"),
		unleash.WithInstanceId("FxiPv5wmNedeSiao1Hh2"),
		unleash.WithAppName("Testing"), // Set to the running environment of your application
		unleash.WithListener(&metricsInterface{}),
	)
	if err != nil {
		logger.Panic("Error init from github: " + err.Error())
		panic(err)
	}
	if unleash.IsEnabled("debug_mode") {
		logger.Info("Debug_mode is enabled")
	}
}

func Get_fflag() {
	config := flagsmith.DefaultConfig()
	client := flagsmith.NewClient("ser.dVQq8WxV2w3iGGbz8DCnHQ", config)
	FFLAG = client
	isEnabled, _ := FFLAG.FeatureEnabled("always_true")
	if !isEnabled {
		logger.Panic("Connection to Flagsmith failed")
		panic("Connection to Flagsmith failed")
	}
}

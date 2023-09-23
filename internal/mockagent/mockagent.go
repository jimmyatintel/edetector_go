package mockagent

import (
	"edetector_go/config"
	"edetector_go/pkg/fflag"
	"edetector_go/pkg/logger"
	"net"
	"os"
)

func init() {
	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Panic("Error loading feature flag")
		panic("Error loading feature flag")
	}
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Panic("Error loading config file: " + err.Error())
		panic(err)
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("WORKER_LOG_FILE"), "server", "SERVER")
		logger.Info("Logger is enabled please check all out info in log file: " + config.Viper.GetString("WORKER_LOG_FILE"))
	}
}

func Main() {
	serverAddr := config.Viper.GetString("WORKING_SERVER_IP") + ":" + config.Viper.GetString("WORKER_DEFAULT_WORKER_PORT")
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		logger.Panic("Error connecting to the server:" + err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	logger.Info("Connected to the server at " + serverAddr)

	// Add your client logic here...
}

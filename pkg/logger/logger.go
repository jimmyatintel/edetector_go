package logger

import (
	"edetector_go/config"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	zapsyslog "github.com/imperfectgo/zap-syslog"
	"github.com/imperfectgo/zap-syslog/syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(path string, hostname string, app string) {
	// file logger
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	fileWriteSyncer := zapcore.AddSync(file)
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)
	// console logger
	stdout := zapcore.AddSync(os.Stdout)
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	// system logger
	syslogCfg := zapsyslog.SyslogEncoderConfig{
		EncoderConfig: zapcore.EncoderConfig{
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		Facility: syslog.LOG_LOCAL0,
		Hostname: hostname,
		PID:      os.Getpid(),
		App:      app,
	}
	syslogEncoder := zapsyslog.NewSyslogEncoder(syslogCfg)
	graylogPath := fmt.Sprintf("%s:%s", config.Viper.GetString("GRAYLOG_HOST"), config.Viper.GetString("GRAYLOG_SYSLOG_PORT"))
	sync, err := zapsyslog.NewConnSyncer("tcp", graylogPath)
	if err != nil {
		fmt.Println(err)
	}
	// set core
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, fileWriteSyncer, zapcore.DebugLevel),
		zapcore.NewCore(syslogEncoder, zapcore.Lock(sync), zapcore.DebugLevel),
	)
	Log = zap.New(core)
}

func Debug(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Error(message, fields...)
}

func Panic(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Panic(message, fields...)
}

func getCallerInfoForLog() (callerFields []zap.Field) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName)

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}

func GinLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		requestURI := c.Request.RequestURI
		Log.Info("GinLog",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("ip", clientIP),
			zap.String("uri", requestURI),
			zap.String("latency", latencyTime.String()),
		)
	}
}

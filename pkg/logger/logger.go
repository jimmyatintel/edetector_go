package logger

import (
	"edetector_go/config"
	"edetector_go/pkg/mariadb"
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
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger
var Service string

func InitLogger(path string, hostname string, app string) {
	// file logger
	// file, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	// fileWriteSyncer := zapcore.AddSync(file)
	rotateCfg := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    5, // megabytes
		MaxBackups: 10,
		MaxAge:     30, // days
	})
	rotatedFileWriteSyncer := zapcore.AddSync(rotateCfg)
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
		zapcore.NewCore(consoleEncoder, stdout, zapcore.InfoLevel),
		zapcore.NewCore(fileEncoder, rotatedFileWriteSyncer, zapcore.DebugLevel),
		zapcore.NewCore(syslogEncoder, zapcore.Lock(sync), zapcore.InfoLevel),
	)
	Log = zap.New(core)
	Service = app
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
	StoreLogToDB("ERROR", message)
}

func Panic(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Panic(message, fields...)
	StoreLogToDB("PANIC", message)
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
		ginInfo := fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s %s", endTime.Format("2006/01/02 - 15:04:05"), statusCode, latencyTime, clientIP, method, requestURI)
		Info(ginInfo)
	}
}

func StoreLogToDB(level string, message string) {
	if len(message) > 250 {
		message = message[:250]
	}
	query := "INSERT INTO log (level, service, content, timestamp) VALUES (?, ?, ?, CURRENT_TIMESTAMP)"
	_, err := mariadb.DB.Exec(query, level, Service, message)
	if err != nil {
		Log.Error("Error storing log to database: " + err.Error())
	}
}

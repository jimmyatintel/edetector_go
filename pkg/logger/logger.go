package logger

import (
	"os"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(path string) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	file, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 644)
	fileWriteSyncer := zapcore.AddSync(file)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel), // you can change the output level here
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	Log = zap.New(core)
}

func Info(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Error(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Log.Warn(message, fields...)
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

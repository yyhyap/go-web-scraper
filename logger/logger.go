package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger = InitializeLogger()

// https://codewithmukesh.com/blog/structured-logging-in-golang-with-zap/
func InitializeLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	// specify the path of log file here
	logFilePath := filepath.Join(".", "log_file")
	err := os.MkdirAll(logFilePath, os.ModePerm)
	if err != nil {
		log.Println("unable to create log file folder")
	}
	logFile, _ := os.OpenFile(fmt.Sprintf("%v/web_scraper_%v.log", logFilePath, time.Now().Format("2006-01-02")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	// set log level here
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, zapcore.WarnLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

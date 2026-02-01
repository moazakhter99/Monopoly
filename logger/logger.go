package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/natefinch/lumberjack"
)

var ZapLogger *zap.SugaredLogger


func Logger() {

	logWrtiter := []zapcore.WriteSyncer{} 

	// Custom time encoder with microseconds and slash format
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006/01/02 15:04:05.000000"))
	}
	
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder ,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// encoder := zapcore.NewConsoleEncoder(encoderConfig)
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleWriter := zapcore.AddSync(os.Stdout)
	logWrtiter = append(logWrtiter, consoleWriter)

	// Getting log file path from the config
	logFilePath := viper.GetString("LOG_FILEPATH")
	lumberjackLogger := &lumberjack.Logger{		// Log file absolute path, os agnostic
		Filename:   filepath.ToSlash(logFilePath),
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	
	log.Printf("Logger Config\n")
	log.Printf("Logger LogFile: %s", lumberjackLogger.Filename)
	log.Printf("| Logger MaxSize: %d MB", lumberjackLogger.MaxSize)
	log.Printf("| Logger MaxBackups: %d", lumberjackLogger.MaxBackups)
	log.Printf("| Logger MaxAge: %d days", lumberjackLogger.MaxAge)
	log.Printf("| Logger Compress: %v\n", lumberjackLogger.Compress)

	fileWriter := zapcore.AddSync(lumberjackLogger)
	logWrtiter = append(logWrtiter, fileWriter)
		
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(logWrtiter...)
	core := zapcore.NewCore(encoder, multiWriteSyncer, zap.DebugLevel)

	ZapLogger = zap.New(zapcore.NewTee(
		core,
	), zap.AddCaller()).Sugar()
	defer ZapLogger.Sync()



}
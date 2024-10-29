package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// Init 初始化全局 logger
func Init() {
	once.Do(func() {
		err := os.MkdirAll("log", 0666)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}
		// 创建日志文件
		infoLog, err := os.OpenFile("log/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		errLog, err := os.OpenFile("log/err.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		// 创建编码器
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
			EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 创建核心
		infoCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(infoLog),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl == zapcore.InfoLevel
			}),
		)

		errorCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(errLog),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		)

		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)

		// 创建 Logger
		logger = zap.New(zapcore.NewTee(infoCore, errorCore, consoleCore), zap.AddCaller())
	})
	fmt.Println("successfully init the logger ")
}

// GetLogger 返回全局 logger 实例
func GetLogger() *zap.Logger {
	if logger == nil {
		Init()
	}
	return logger
}

// Sync 刷新缓冲日志条目
func Sync() {
	if logger != nil {
		logger.Sync()
	}
}

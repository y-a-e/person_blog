package core

import (
	"log"
	"os"
	"server/global"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 初InitLogger 初始化并返回一个基于配置设置的新 zap.Logger 实例
func InitLogger() *zap.Logger {
	// 直接引用位置文件数据
	zapCfg := global.Config.Zap

	// 调用函数将日志写入文件
	writeSyaner := getLogWriter(zapCfg.Filename, zapCfg.MaxSize, zapCfg.MaxBackups, zapCfg.MaxAge)
	// 控制台是否需要打印
	if zapCfg.IsConsolePrint {
		writeSyaner = zapcore.NewMultiWriteSyncer(writeSyaner, zapcore.AddSync(os.Stdout))
	}

	// 创建日志格式化的编码器
	encoder := getEncoder()

	// 根据配置确定日志级别
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(zapCfg.Level))
	if err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}

	// 创建对象封装
	core := zapcore.NewCore(encoder, writeSyaner, logLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

// 引入lumberJack包，将日志文件切割归档
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 配置格式写死
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

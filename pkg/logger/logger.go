package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// Init 初始化日志
func Init(level string, format string) error {
	var config zap.Config

	// 根据日志级别设置配置
	switch level {
	case "debug":
		config = zap.NewDevelopmentConfig()
	case "info", "warn", "error":
		config = zap.NewProductionConfig()
	default:
		config = zap.NewProductionConfig()
	}

	// 设置日志级别
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}

	// 设置输出格式
	if format == "text" {
		config.Encoding = "console"
	} else {
		config.Encoding = "json"
	}

	// 创建日志记录器
	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("创建日志记录器失败: %w", err)
	}

	globalLogger = logger
	return nil
}

// GetLogger 获取全局日志记录器
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// 如果未初始化，创建一个默认的日志记录器
		logger, _ := zap.NewProduction()
		globalLogger = logger
	}
	return globalLogger
}

// Debug 记录调试日志
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info 记录信息日志
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn 记录警告日志
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error 记录错误日志
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal 记录致命错误日志并退出
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Sync 同步日志缓冲区
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Close 关闭日志记录器
func Close() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Printf 格式化打印日志
func Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
}

// Println 打印日志行
func Println(args ...interface{}) {
	fmt.Fprintln(os.Stdout, args...)
}

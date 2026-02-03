package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInit_DebugLevel(t *testing.T) {
	// 初始化调试级别日志
	err := Init("debug", "json")
	require.NoError(t, err)

	// 验证日志记录器已初始化
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestInit_InfoLevel(t *testing.T) {
	// 初始化信息级别日志
	err := Init("info", "json")
	require.NoError(t, err)

	// 验证日志记录器已初始化
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestInit_WarnLevel(t *testing.T) {
	// 初始化警告级别日志
	err := Init("warn", "json")
	require.NoError(t, err)

	// 验证日志记录器已初始化
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestInit_ErrorLevel(t *testing.T) {
	// 初始化错误级别日志
	err := Init("error", "json")
	require.NoError(t, err)

	// 验证日志记录器已初始化
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestInit_TextFormat(t *testing.T) {
	// 初始化文本格式日志
	err := Init("info", "text")
	require.NoError(t, err)

	// 验证日志记录器已初始化
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestInit_JSONFormat(t *testing.T) {
	// 初始化 JSON 格式日志
	err := Init("info", "json")
	require.NoError(t, err)

	// 验证日志记录器已初始化
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestGetLogger_DefaultLogger(t *testing.T) {
	// 重置全局日志记录器
	globalLogger = nil

	// 获取日志记录器
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestGetLogger_Consistency(t *testing.T) {
	// 初始化日志
	err := Init("info", "json")
	require.NoError(t, err)

	// 多次获取日志记录器应该返回相同的实例
	logger1 := GetLogger()
	logger2 := GetLogger()
	assert.Equal(t, logger1, logger2)
}

func TestDebug(t *testing.T) {
	// 初始化日志
	err := Init("debug", "json")
	require.NoError(t, err)

	// 测试 Debug 日志
	Debug("test debug message", zap.String("key", "value"))
	// 验证没有 panic
}

func TestInfo(t *testing.T) {
	// 初始化日志
	err := Init("info", "json")
	require.NoError(t, err)

	// 测试 Info 日志
	Info("test info message", zap.String("key", "value"))
	// 验证没有 panic
}

func TestWarn(t *testing.T) {
	// 初始化日志
	err := Init("warn", "json")
	require.NoError(t, err)

	// 测试 Warn 日志
	Warn("test warn message", zap.String("key", "value"))
	// 验证没有 panic
}

func TestError(t *testing.T) {
	// 初始化日志
	err := Init("error", "json")
	require.NoError(t, err)

	// 测试 Error 日志
	Error("test error message", zap.String("key", "value"))
	// 验证没有 panic
}

func TestSync(t *testing.T) {
	// 初始化日志
	err := Init("info", "json")
	require.NoError(t, err)

	// 同步日志缓冲区
	err = Sync()
	assert.NoError(t, err)
}

func TestClose(t *testing.T) {
	// 初始化日志
	err := Init("info", "json")
	require.NoError(t, err)

	// 关闭日志记录器
	err = Close()
	assert.NoError(t, err)
}

func TestSync_WithoutInit(t *testing.T) {
	// 重置全局日志记录器
	globalLogger = nil

	// 同步日志缓冲区（未初始化）
	err := Sync()
	assert.NoError(t, err)
}

func TestClose_WithoutInit(t *testing.T) {
	// 重置全局日志记录器
	globalLogger = nil

	// 关闭日志记录器（未初始化）
	err := Close()
	assert.NoError(t, err)
}

func TestPrintf(t *testing.T) {
	// 测试 Printf
	Printf("test printf: %s\n", "message")
	// 验证没有 panic
}

func TestPrintln(t *testing.T) {
	// 测试 Println
	Println("test println message")
	// 验证没有 panic
}

func TestInit_InvalidLevel(t *testing.T) {
	// 初始化无效的日志级别（应该使用默认值）
	err := Init("invalid", "json")
	require.NoError(t, err)

	// 验证日志记录器已初始化
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestMultipleInit(t *testing.T) {
	// 多次初始化日志
	err1 := Init("info", "json")
	require.NoError(t, err1)

	err2 := Init("debug", "text")
	require.NoError(t, err2)

	// 验证日志记录器已更新
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestLoggerWithFields(t *testing.T) {
	// 初始化日志
	err := Init("info", "json")
	require.NoError(t, err)

	// 测试带字段的日志
	Info("test message",
		zap.String("user", "john"),
		zap.Int("age", 30),
		zap.Bool("active", true),
	)
	// 验证没有 panic
}

func TestLoggerMultipleCalls(t *testing.T) {
	// 初始化日志
	err := Init("info", "json")
	require.NoError(t, err)

	// 多次调用日志方法
	for i := 0; i < 10; i++ {
		Info("test message", zap.Int("iteration", i))
	}
	// 验证没有 panic
}

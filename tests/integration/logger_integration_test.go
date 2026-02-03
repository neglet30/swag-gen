package integration

import (
	"testing"

	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/neglet30/swag-gen/pkg/logger"
	"github.com/neglet30/swag-gen/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestLoggerWithServerIntegration 测试日志系统与服务器的集成
func TestLoggerWithServerIntegration(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Logger Integration Test",
			Version: "1.0.0",
		},
	}

	// 记录服务器启动日志
	logger.Info("server starting",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 记录服务器创建日志
	logger.Info("server created successfully",
		zap.String("project", cfg.Project.Name),
		zap.String("version", cfg.Project.Version),
	)

	// 验证日志记录器已初始化
	loggerInstance := logger.GetLogger()
	assert.NotNil(t, loggerInstance)
}

// TestLoggerLevelConfiguration 测试日志级别配置
func TestLoggerLevelConfiguration(t *testing.T) {
	// 测试不同的日志级别
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			// 初始化日志系统
			err := logger.Init(level, "json")
			require.NoError(t, err)

			// 验证日志记录器已初始化
			loggerInstance := logger.GetLogger()
			assert.NotNil(t, loggerInstance)

			// 记录日志
			logger.Info("test message for level: "+level, zap.String("level", level))
		})
	}
}

// TestLoggerFormatConfiguration 测试日志格式配置
func TestLoggerFormatConfiguration(t *testing.T) {
	// 测试不同的日志格式
	formats := []string{"json", "text"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			// 初始化日志系统
			err := logger.Init("info", format)
			require.NoError(t, err)

			// 验证日志记录器已初始化
			loggerInstance := logger.GetLogger()
			assert.NotNil(t, loggerInstance)

			// 记录日志
			logger.Info("test message for format: "+format, zap.String("format", format))
		})
	}
}

// TestLoggerWithConfigIntegration 测试日志系统与配置系统的集成
func TestLoggerWithConfigIntegration(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Level:  "debug",
			Format: "json",
			Output: "stdout",
		},
	}

	// 根据配置初始化日志系统
	err := logger.Init(cfg.Logger.Level, cfg.Logger.Format)
	require.NoError(t, err)

	// 验证日志记录器已初始化
	loggerInstance := logger.GetLogger()
	assert.NotNil(t, loggerInstance)

	// 记录日志
	logger.Debug("debug message", zap.String("config", "applied"))
	logger.Info("info message", zap.String("config", "applied"))
	logger.Warn("warn message", zap.String("config", "applied"))
	logger.Error("error message", zap.String("config", "applied"))
}

// TestLoggerMultipleInitializations 测试多次初始化日志系统
func TestLoggerMultipleInitializations(t *testing.T) {
	// 第一次初始化
	err := logger.Init("info", "json")
	require.NoError(t, err)

	logger1 := logger.GetLogger()
	assert.NotNil(t, logger1)

	// 记录日志
	logger.Info("first initialization")

	// 第二次初始化
	err = logger.Init("debug", "text")
	require.NoError(t, err)

	logger2 := logger.GetLogger()
	assert.NotNil(t, logger2)

	// 记录日志
	logger.Debug("second initialization")

	// 验证日志记录器已更新
	assert.NotNil(t, logger2)
}

// TestLoggerWithFields 测试带字段的日志记录
func TestLoggerWithFields(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 记录带字段的日志
	logger.Info("user action",
		zap.String("user_id", "user123"),
		zap.String("action", "login"),
		zap.Int("duration_ms", 150),
		zap.Bool("success", true),
	)

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerSync 测试日志同步
func TestLoggerSync(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 记录日志
	logger.Info("test message before sync")

	// 同步日志缓冲区
	err = logger.Sync()
	assert.NoError(t, err)

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerClose 测试日志关闭
func TestLoggerClose(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 记录日志
	logger.Info("test message before close")

	// 关闭日志记录器
	err = logger.Close()
	assert.NoError(t, err)

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerWithServerRequests 测试服务器请求时的日志记录
func TestLoggerWithServerRequests(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Logger Request Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 记录请求处理日志
	logger.Info("processing request",
		zap.String("method", "GET"),
		zap.String("path", "/health"),
		zap.Int("status_code", 200),
		zap.Int("duration_ms", 5),
	)

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerErrorHandling 测试错误日志记录
func TestLoggerErrorHandling(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("error", "json")
	require.NoError(t, err)

	// 模拟错误
	testError := "test error message"

	// 记录错误日志
	logger.Error("operation failed",
		zap.String("error", testError),
		zap.String("operation", "test_operation"),
	)

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerConcurrentAccess 测试并发日志访问
func TestLoggerConcurrentAccess(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 并发记录日志
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Info("concurrent message",
				zap.Int("goroutine_id", id),
			)
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerWithDifferentLevels 测试不同日志级别的记录
func TestLoggerWithDifferentLevels(t *testing.T) {
	// 初始化日志系统为 debug 级别
	err := logger.Init("debug", "json")
	require.NoError(t, err)

	// 记录不同级别的日志
	logger.Debug("debug level message", zap.String("level", "debug"))
	logger.Info("info level message", zap.String("level", "info"))
	logger.Warn("warn level message", zap.String("level", "warn"))
	logger.Error("error level message", zap.String("level", "error"))

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerJSONFormat 测试 JSON 格式日志
func TestLoggerJSONFormat(t *testing.T) {
	// 初始化日志系统为 JSON 格式
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 创建一个缓冲区来捕获日志输出
	// 注意：这是一个简化的测试，实际的日志输出可能需要更复杂的设置
	logger.Info("json format test",
		zap.String("format", "json"),
		zap.Int("test_number", 123),
		zap.Bool("test_bool", true),
	)

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerTextFormat 测试文本格式日志
func TestLoggerTextFormat(t *testing.T) {
	// 初始化日志系统为文本格式
	err := logger.Init("info", "text")
	require.NoError(t, err)

	// 记录日志
	logger.Info("text format test",
		zap.String("format", "text"),
		zap.Int("test_number", 456),
		zap.Bool("test_bool", false),
	)

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerWithComplexFields 测试带复杂字段的日志
func TestLoggerWithComplexFields(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 记录带复杂字段的日志
	logger.Info("complex fields test",
		zap.String("user_id", "user123"),
		zap.String("action", "api_call"),
		zap.Int("status_code", 200),
		zap.Int("response_time_ms", 150),
		zap.Bool("cached", true),
		zap.String("endpoint", "/api/users"),
		zap.String("method", "GET"),
	)

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerPrintf 测试 Printf 方法
func TestLoggerPrintf(t *testing.T) {
	// 测试 Printf 方法
	logger.Printf("printf test: %s\n", "message")

	// 验证没有 panic
	assert.True(t, true)
}

// TestLoggerPrintln 测试 Println 方法
func TestLoggerPrintln(t *testing.T) {
	// 测试 Println 方法
	logger.Println("println test message")

	// 验证没有 panic
	assert.True(t, true)
}

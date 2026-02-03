package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/neglet30/swag-gen/pkg/logger"
	"github.com/neglet30/swag-gen/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func init() {
	// 初始化日志
	logger.Init("info", "json")
}

// TestConfigLoadAndServerInit 测试配置加载与服务器初始化集成
func TestConfigLoadAndServerInit(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: "127.0.0.1"
  port: 9090
  env: "development"
  read_timeout: 30
  write_timeout: 30

project:
  name: "Integration Test API"
  version: "1.0.0"
  description: "Integration Test Description"
  base_path: "/api/v1"

parser:
  enable_cache: true
  cache_ttl: 3600
  max_concurrent: 4
  exclude_dirs:
    - "vendor"
    - "test"

swagger:
  version: "3.0.0"
  default_version: "1.0.0"
  default_title: "Integration Test API"

logger:
  level: "info"
  format: "json"
  output: "stdout"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// 加载配置
	cfg, err := config.Load(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证配置值
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "Integration Test API", cfg.Project.Name)
	assert.Equal(t, "1.0.0", cfg.Project.Version)

	// 使用配置创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 验证服务器配置正确
	assert.Equal(t, cfg, srv.GetConfig())
}

// TestEnvironmentVariableOverride 测试环境变量覆盖配置
func TestEnvironmentVariableOverride(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: "0.0.0.0"
  port: 8080
  env: "development"

project:
  name: "Test API"
  version: "1.0.0"

logger:
  level: "info"
  format: "json"
  output: "stdout"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// 设置环境变量
	os.Setenv("SWAG_GEN_PORT", "7070")
	os.Setenv("SWAG_GEN_HOST", "192.168.1.1")
	os.Setenv("SWAG_GEN_LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("SWAG_GEN_PORT")
		os.Unsetenv("SWAG_GEN_HOST")
		os.Unsetenv("SWAG_GEN_LOG_LEVEL")
	}()

	// 加载配置
	cfg, err := config.Load(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证环境变量覆盖
	assert.Equal(t, "192.168.1.1", cfg.Server.Host)
	assert.Equal(t, 7070, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Logger.Level)

	// 验证其他配置值保持不变
	assert.Equal(t, "Test API", cfg.Project.Name)
	assert.Equal(t, "1.0.0", cfg.Project.Version)
}

// TestConfigWithLoggerIntegration 测试配置与日志系统集成
func TestConfigWithLoggerIntegration(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
logger:
  level: "debug"
  format: "json"
  output: "stdout"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// 加载配置
	cfg, err := config.Load(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证日志配置
	assert.Equal(t, "debug", cfg.Logger.Level)
	assert.Equal(t, "json", cfg.Logger.Format)
	assert.Equal(t, "stdout", cfg.Logger.Output)

	// 根据配置初始化日志系统
	err = logger.Init(cfg.Logger.Level, cfg.Logger.Format)
	require.NoError(t, err)

	// 验证日志记录器已初始化
	loggerInstance := logger.GetLogger()
	assert.NotNil(t, loggerInstance)

	// 记录日志
	logger.Info("test message from config")
}

// TestConfigWithServerAndLogger 测试配置、服务器和日志系统的完整集成
func TestConfigWithServerAndLogger(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: "127.0.0.1"
  port: 8888
  env: "development"

project:
  name: "Full Integration Test API"
  version: "1.0.0"
  description: "Full Integration Test"

logger:
  level: "info"
  format: "json"
  output: "stdout"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// 加载配置
	cfg, err := config.Load(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// 初始化日志系统
	err = logger.Init(cfg.Logger.Level, cfg.Logger.Format)
	require.NoError(t, err)

	// 记录启动日志
	logger.Info("starting server",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 验证服务器配置
	assert.Equal(t, cfg, srv.GetConfig())
	assert.Equal(t, "Full Integration Test API", srv.GetConfig().Project.Name)
	assert.Equal(t, 8888, srv.GetConfig().Server.Port)

	// 记录完成日志
	logger.Info("server created successfully")
}

// TestMultipleConfigLoads 测试多次加载配置
func TestMultipleConfigLoads(t *testing.T) {
	// 创建第一个配置文件
	tmpDir := t.TempDir()
	configPath1 := filepath.Join(tmpDir, "config1.yaml")

	configContent1 := `
server:
  port: 8080

project:
  name: "API 1"
  version: "1.0.0"
`

	err := os.WriteFile(configPath1, []byte(configContent1), 0644)
	require.NoError(t, err)

	// 加载第一个配置
	cfg1, err := config.Load(configPath1)
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg1.Server.Port)
	assert.Equal(t, "API 1", cfg1.Project.Name)

	// 创建第二个配置文件
	configPath2 := filepath.Join(tmpDir, "config2.yaml")

	configContent2 := `
server:
  port: 9090

project:
  name: "API 2"
  version: "2.0.0"
`

	err = os.WriteFile(configPath2, []byte(configContent2), 0644)
	require.NoError(t, err)

	// 加载第二个配置
	cfg2, err := config.Load(configPath2)
	require.NoError(t, err)
	assert.Equal(t, 9090, cfg2.Server.Port)
	assert.Equal(t, "API 2", cfg2.Project.Name)

	// 验证两个配置是独立的
	assert.NotEqual(t, cfg1.Server.Port, cfg2.Server.Port)
	assert.NotEqual(t, cfg1.Project.Name, cfg2.Project.Name)
}

// TestConfigDefaultValues 测试配置默认值
func TestConfigDefaultValues(t *testing.T) {
	// 加载空配置（使用默认值）
	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证所有默认值都已设置
	assert.NotEmpty(t, cfg.Server.Host)
	assert.Greater(t, cfg.Server.Port, 0)
	assert.NotEmpty(t, cfg.Server.Env)
	assert.Greater(t, cfg.Server.ReadTimeout, 0)
	assert.Greater(t, cfg.Server.WriteTimeout, 0)

	assert.NotEmpty(t, cfg.Project.Name)
	assert.NotEmpty(t, cfg.Project.Version)
	assert.NotEmpty(t, cfg.Project.BasePath)

	assert.NotEmpty(t, cfg.Logger.Level)
	assert.NotEmpty(t, cfg.Logger.Format)
	assert.NotEmpty(t, cfg.Logger.Output)

	// 使用默认配置创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)
}

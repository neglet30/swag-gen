package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/neglet30/swag-gen/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// 初始化日志
	logger.Init("info", "json")
}

func TestNew(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)
	assert.NotNil(t, server)
	assert.NotNil(t, server.engine)
	assert.Equal(t, cfg, server.config)
}

func TestHealthHandler(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])
	assert.NotNil(t, response["data"])
}

func TestGetSwaggerHandler(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:        "Test API",
			Version:     "1.0.0",
			Description: "Test Description",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/swagger", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "3.0.0", data["openapi"])

	info := data["info"].(map[string]interface{})
	assert.Equal(t, "Test API", info["title"])
	assert.Equal(t, "1.0.0", info["version"])
	assert.Equal(t, "Test Description", info["description"])
}

func TestGetEndpointsHandler(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/api/endpoints", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["total"])
	assert.NotNil(t, data["items"])
}

func TestTestAPIHandler(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("POST", "/api/test", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "test-001", data["id"])
	assert.Equal(t, float64(200), data["statusCode"])
}

func TestGetTestHistoryHandler(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/api/test/history", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["total"])
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(10), data["pageSize"])
	assert.NotNil(t, data["items"])
}

func TestGetTestDetailHandler(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/api/test/test-123", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "test-123", data["id"])
}

func TestClearTestHistoryHandler(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("DELETE", "/api/test/history", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])
	assert.Nil(t, response["data"])
}

func TestCORSHeaders(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证 CORS 头
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

func TestOPTIONSRequest(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 创建 OPTIONS 请求
	req, err := http.NewRequest("OPTIONS", "/health", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	server.engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestProductionMode(t *testing.T) {
	// 创建生产环境配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "production",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)
	assert.NotNil(t, server)
	assert.NotNil(t, server.engine)
}

func TestDevelopmentMode(t *testing.T) {
	// 创建开发环境配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)
	assert.NotNil(t, server)
	assert.NotNil(t, server.engine)
}

func TestMultipleRequests(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Env:  "development",
		},
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	server := New(cfg)

	// 执行多个请求
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/neglet30/swag-gen/pkg/logger"
	"github.com/neglet30/swag-gen/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompleteRequestFlow 测试完整的请求处理流程
func TestCompleteRequestFlow(t *testing.T) {
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
			Name:    "Complete Flow Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 发送请求
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

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

// TestConcurrentRequests 测试并发请求处理
func TestConcurrentRequests(t *testing.T) {
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
			Name:    "Concurrent Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 并发发送请求
	numRequests := 50
	var wg sync.WaitGroup
	results := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req, err := http.NewRequest("GET", "/health", nil)
			if err != nil {
				results <- false
				return
			}

			w := httptest.NewRecorder()
			srv.GetEngine().ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				results <- true
			} else {
				results <- false
			}
		}()
	}

	// 等待所有请求完成
	wg.Wait()
	close(results)

	// 验证所有请求都成功
	successCount := 0
	for result := range results {
		if result {
			successCount++
		}
	}

	assert.Equal(t, numRequests, successCount)
}

// TestDifferentHTTPMethods 测试不同的 HTTP 方法
func TestDifferentHTTPMethods(t *testing.T) {
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
			Name:    "HTTP Methods Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 测试 GET 请求
	t.Run("GET", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		srv.GetEngine().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 测试 POST 请求
	t.Run("POST", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		srv.GetEngine().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 测试 DELETE 请求
	t.Run("DELETE", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/api/test/history", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		srv.GetEngine().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 测试 OPTIONS 请求
	t.Run("OPTIONS", func(t *testing.T) {
		req, err := http.NewRequest("OPTIONS", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		srv.GetEngine().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

// TestCORSMiddlewareIntegration 测试 CORS 中间件集成
func TestCORSMiddlewareIntegration(t *testing.T) {
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
			Name:    "CORS Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 发送带有 Origin 头的请求
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://localhost:3000")

	// 执行请求
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	// 验证 CORS 头
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

// TestPreflight 测试预检请求
func TestPreflight(t *testing.T) {
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
			Name:    "Preflight Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 发送 OPTIONS 预检请求
	req, err := http.NewRequest("OPTIONS", "/api/test", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	// 执行请求
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestLoggerMiddlewareIntegration 测试日志中间件集成
func TestLoggerMiddlewareIntegration(t *testing.T) {
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
			Name:    "Logger Middleware Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 发送请求
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证日志中间件已应用（通过检查响应头）
	assert.NotEmpty(t, w.Header().Get("Content-Type"))
}

// TestMultipleEndpoints 测试多个端点
func TestMultipleEndpoints(t *testing.T) {
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
			Name:    "Multiple Endpoints Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 测试多个端点
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/health"},
		{"GET", "/swagger"},
		{"GET", "/api/endpoints"},
		{"POST", "/api/test"},
		{"GET", "/api/test/history"},
		{"GET", "/api/test/test-123"},
		{"DELETE", "/api/test/history"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			req, err := http.NewRequest(endpoint.method, endpoint.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			srv.GetEngine().ServeHTTP(w, req)

			// 验证响应状态码
			assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNoContent)
		})
	}
}

// TestResponseFormat 测试响应格式
func TestResponseFormat(t *testing.T) {
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
			Name:    "Response Format Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 发送请求
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	// 验证响应格式
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// 验证响应结构
	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
}

// TestErrorResponse 测试错误响应
func TestErrorResponse(t *testing.T) {
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
			Name:    "Error Response Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 发送请求到不存在的端点
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestRequestWithBody 测试带有请求体的请求
func TestRequestWithBody(t *testing.T) {
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
			Name:    "Request Body Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 创建请求体
	body := []byte(`{"method":"GET","url":"http://localhost:8000/api/users"}`)

	// 发送 POST 请求
	req, err := http.NewRequest("POST", "/api/test", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestRequestHeaders 测试请求头
func TestRequestHeaders(t *testing.T) {
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
			Name:    "Request Headers Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 发送带有自定义头的请求
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("X-Custom-Header", "custom-value")

	// 执行请求
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestResponseTime 测试响应时间
func TestResponseTime(t *testing.T) {
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
			Name:    "Response Time Test",
			Version: "1.0.0",
		},
	}

	// 创建服务器
	srv := server.New(cfg)
	assert.NotNil(t, srv)

	// 发送请求并测量响应时间
	start := time.Now()

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	duration := time.Since(start)

	// 验证响应时间在合理范围内（< 100ms）
	assert.Less(t, duration, 100*time.Millisecond)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestServerWithDifferentEnvironments 测试不同环境的服务器
func TestServerWithDifferentEnvironments(t *testing.T) {
	// 初始化日志系统
	err := logger.Init("info", "json")
	require.NoError(t, err)

	// 测试开发环境
	t.Run("development", func(t *testing.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
				Port: 8080,
				Env:  "development",
			},
			Project: config.ProjectConfig{
				Name:    "Dev Test",
				Version: "1.0.0",
			},
		}

		srv := server.New(cfg)
		assert.NotNil(t, srv)

		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		srv.GetEngine().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 测试生产环境
	t.Run("production", func(t *testing.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
				Port: 8080,
				Env:  "production",
			},
			Project: config.ProjectConfig{
				Name:    "Prod Test",
				Version: "1.0.0",
			},
		}

		srv := server.New(cfg)
		assert.NotNil(t, srv)

		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		srv.GetEngine().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/neglet30/swag-gen/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// 初始化日志
	logger.Init("info", "json")
}

func TestLoggerMiddleware(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(LoggerMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggerMiddleware_DifferentMethods(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(LoggerMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.PUT("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.DELETE("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试不同的 HTTP 方法
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	for _, method := range methods {
		req, err := http.NewRequest(method, "/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestLoggerMiddleware_DifferentPaths(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(LoggerMiddleware())

	// 添加测试路由
	engine.GET("/api/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.GET("/api/posts", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试不同的路径
	paths := []string{"/api/users", "/api/posts"}
	for _, path := range paths {
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestLoggerMiddleware_StatusCodes(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(LoggerMiddleware())

	// 添加测试路由
	engine.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.GET("/created", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "created"})
	})
	engine.GET("/bad-request", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})
	engine.GET("/not-found", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})
	engine.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	})

	// 测试不同的状态码
	tests := []struct {
		path       string
		statusCode int
	}{
		{"/ok", http.StatusOK},
		{"/created", http.StatusCreated},
		{"/bad-request", http.StatusBadRequest},
		{"/not-found", http.StatusNotFound},
		{"/error", http.StatusInternalServerError},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", test.path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, test.statusCode, w.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(CORSMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证 CORS 头
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

func TestCORSMiddleware_AllowedMethods(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(CORSMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.PUT("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.DELETE("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 测试不同的 HTTP 方法
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	for _, method := range methods {
		req, err := http.NewRequest(method, "/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		// 验证 CORS 头
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSMiddleware_OPTIONSRequest(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(CORSMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建 OPTIONS 请求
	req, err := http.NewRequest("OPTIONS", "/test", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 验证 CORS 头
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_AllowedHeaders(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(CORSMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证允许的请求头
	allowedHeaders := w.Header().Get("Access-Control-Allow-Headers")
	assert.Contains(t, allowedHeaders, "Content-Type")
	assert.Contains(t, allowedHeaders, "Authorization")
}

func TestErrorHandlerMiddleware(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(ErrorHandlerMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestErrorHandlerMiddleware_WithError(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(ErrorHandlerMiddleware())

	// 添加测试路由
	engine.GET("/error", func(c *gin.Context) {
		err := gin.Error{Err: nil, Type: gin.ErrorTypeBind}
		c.Error(&err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})

	// 创建测试请求
	req, err := http.NewRequest("GET", "/error", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMultipleMiddlewares(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(LoggerMiddleware())
	engine.Use(CORSMiddleware())
	engine.Use(ErrorHandlerMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证 CORS 头
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestMiddlewareOrder(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()

	// 添加中间件（顺序很重要）
	engine.Use(LoggerMiddleware())
	engine.Use(CORSMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddleware_CustomOrigin(t *testing.T) {
	// 创建 Gin 引擎
	engine := gin.New()
	engine.Use(CORSMiddleware())

	// 添加测试路由
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 创建带有 Origin 头的测试请求
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://example.com")

	// 执行请求
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 验证 CORS 头（应该允许所有源）
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

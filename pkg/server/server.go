package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/neglet30/swag-gen/pkg/logger"
)

// Server HTTP 服务器
type Server struct {
	engine *gin.Engine
	config *config.Config
}

// New 创建新的服务器实例
func New(cfg *config.Config) *Server {
	// 设置 Gin 模式
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

	// 添加中间件
	engine.Use(gin.Recovery())
	engine.Use(LoggerMiddleware())
	engine.Use(CORSMiddleware())

	server := &Server{
		engine: engine,
		config: cfg,
	}

	// 注册路由
	server.registerRoutes()

	return server
}

// registerRoutes 注册路由
func (s *Server) registerRoutes() {
	// 健康检查
	s.engine.GET("/health", s.healthHandler)

	// Swagger 文档 API
	s.engine.GET("/swagger", s.getSwaggerHandler)
	s.engine.GET("/swagger/ui", s.getSwaggerUIHandler)
	s.engine.GET("/api/endpoints", s.getEndpointsHandler)

	// API 测试 API
	s.engine.POST("/api/test", s.testAPIHandler)
	s.engine.GET("/api/test/history", s.getTestHistoryHandler)
	s.engine.GET("/api/test/:testId", s.getTestDetailHandler)
	s.engine.DELETE("/api/test/history", s.clearTestHistoryHandler)
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	logger.Info(fmt.Sprintf("启动服务器: %s", addr))
	return s.engine.Run(addr)
}

// Stop 停止服务器
func (s *Server) Stop() error {
	logger.Info("停止服务器")
	return nil
}

// GetEngine 获取 Gin 引擎
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// GetConfig 获取服务器配置
func (s *Server) GetConfig() *config.Config {
	return s.config
}

// 处理器方法

// healthHandler 健康检查处理器
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"status":  "healthy",
			"version": "0.1.0",
		},
	})
}

// getSwaggerHandler 获取 Swagger 文档处理器
func (s *Server) getSwaggerHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"openapi": "3.0.0",
			"info": gin.H{
				"title":       s.config.Project.Name,
				"version":     s.config.Project.Version,
				"description": s.config.Project.Description,
			},
			"paths": gin.H{},
		},
	})
}

// getSwaggerUIHandler 获取 Swagger UI 处理器
func (s *Server) getSwaggerUIHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "swagger-ui.html", gin.H{
		"title": s.config.Project.Name,
	})
}

// getEndpointsHandler 获取所有端点处理器
func (s *Server) getEndpointsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total": 0,
			"items": []interface{}{},
		},
	})
}

// testAPIHandler 测试 API 处理器
func (s *Server) testAPIHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":         "test-001",
			"statusCode": 200,
			"duration":   0,
		},
	})
}

// getTestHistoryHandler 获取测试历史处理器
func (s *Server) getTestHistoryHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total":    0,
			"page":     1,
			"pageSize": 10,
			"items":    []interface{}{},
		},
	})
}

// getTestDetailHandler 获取测试详情处理器
func (s *Server) getTestDetailHandler(c *gin.Context) {
	testID := c.Param("testId")
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id": testID,
		},
	})
}

// clearTestHistoryHandler 清空测试历史处理器
func (s *Server) clearTestHistoryHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    nil,
	})
}

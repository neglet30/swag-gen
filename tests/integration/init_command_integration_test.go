package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/neglet30/swag-gen/pkg/logger"
	"github.com/neglet30/swag-gen/pkg/output"
	"github.com/neglet30/swag-gen/pkg/parser"
	"github.com/neglet30/swag-gen/pkg/swagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitCommandIntegration_CompleteFlow 测试完整的初始化流程
func TestInitCommandIntegration_CompleteFlow(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	projectDir, err := os.MkdirTemp("", "init-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// 创建临时输出目录
	outputDir, err := os.MkdirTemp("", "init-output-*")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	// 创建示例 API 文件
	apiFile := filepath.Join(projectDir, "api.go")
	apiContent := `package api

// GetUsers 获取所有用户
// @Router /api/users [GET]
// @Summary 获取所有用户
// @Description 获取系统中的所有用户
// @Tags User
// @Success 200 {array} User
// @Failure 500 {object} ErrorResponse
func GetUsers() {}

// CreateUser 创建用户
// @Router /api/users [POST]
// @Summary 创建用户
// @Description 创建新用户
// @Tags User
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
func CreateUser() {}

// User 用户模型
type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
}
`
	err = os.WriteFile(apiFile, []byte(apiContent), 0644)
	require.NoError(t, err)

	// 创建配置
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:        "Test API",
			Version:     "1.0.0",
			Description: "Test API Description",
		},
		Parser: config.ParserConfig{
			EnableCache:   true,
			CacheTTL:      3600,
			MaxConcurrent: 4,
			ExcludeDirs:   []string{"vendor", "node_modules"},
		},
	}

	// 创建解析器
	p := parser.NewParser(cfg, log)

	// 解析项目
	endpoints, err := p.ParseProject(projectDir)
	require.NoError(t, err)
	assert.Equal(t, 2, len(endpoints), "应该找到 2 个端点")

	// 创建 Swagger 构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "Test API Description")

	// 添加所有端点
	for _, endpoint := range endpoints {
		err := builder.AddEndpoint(endpoint)
		require.NoError(t, err)
	}

	// 构建文档
	doc := builder.Build()
	assert.NotNil(t, doc)

	// 写入输出文件
	writer := output.NewWriter(outputDir)

	// 写入 Swagger 文档
	err = writer.WriteSwagger(doc, "swagger", "json")
	require.NoError(t, err)

	// 验证 Swagger 文件
	swaggerFile := filepath.Join(outputDir, "swagger.json")
	assert.FileExists(t, swaggerFile)

	// 写入配置文件
	outputConfig := output.NewConfig("Test API", "1.0.0", "Test API Description")
	outputConfig.SetParserPath(projectDir)
	outputConfig.SetOutputPath(outputDir)
	outputConfig.SetOutputFormat("json")

	err = writer.WriteConfig(outputConfig, "swag-gen.yaml")
	require.NoError(t, err)

	// 验证配置文件
	configFile := filepath.Join(outputDir, "swag-gen.yaml")
	assert.FileExists(t, configFile)

	// 写入 README
	err = writer.WriteREADME("README.md", "Test API", "Test API Description")
	require.NoError(t, err)

	// 验证 README 文件
	readmeFile := filepath.Join(outputDir, "README.md")
	assert.FileExists(t, readmeFile)
}

// TestInitCommandIntegration_JSONFormat 测试 JSON 格式输出
func TestInitCommandIntegration_JSONFormat(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	projectDir, err := os.MkdirTemp("", "init-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// 创建临时输出目录
	outputDir, err := os.MkdirTemp("", "init-output-*")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	// 创建示例 API 文件
	apiFile := filepath.Join(projectDir, "api.go")
	apiContent := `package api

// GetUsers 获取用户
// @Router /api/users [GET]
// @Summary 获取用户
// @Tags User
// @Success 200 {array} User
func GetUsers() {}

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	err = os.WriteFile(apiFile, []byte(apiContent), 0644)
	require.NoError(t, err)

	// 创建配置
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建解析器
	p := parser.NewParser(cfg, log)

	// 解析项目
	endpoints, err := p.ParseProject(projectDir)
	require.NoError(t, err)

	// 创建 Swagger 构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 添加所有端点
	for _, endpoint := range endpoints {
		err := builder.AddEndpoint(endpoint)
		require.NoError(t, err)
	}

	// 构建文档
	doc := builder.Build()

	// 写入输出文件
	writer := output.NewWriter(outputDir)

	// 写入 JSON 格式
	err = writer.WriteSwagger(doc, "swagger", "json")
	require.NoError(t, err)

	// 验证文件
	swaggerFile := filepath.Join(outputDir, "swagger.json")
	assert.FileExists(t, swaggerFile)

	// 读取文件内容
	content, err := os.ReadFile(swaggerFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "openapi")
	assert.Contains(t, string(content), "Test API")
}

// TestInitCommandIntegration_YAMLFormat 测试 YAML 格式输出
func TestInitCommandIntegration_YAMLFormat(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	projectDir, err := os.MkdirTemp("", "init-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// 创建临时输出目录
	outputDir, err := os.MkdirTemp("", "init-output-*")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	// 创建示例 API 文件
	apiFile := filepath.Join(projectDir, "api.go")
	apiContent := `package api

// GetUsers 获取用户
// @Router /api/users [GET]
// @Summary 获取用户
// @Tags User
// @Success 200 {array} User
func GetUsers() {}

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	err = os.WriteFile(apiFile, []byte(apiContent), 0644)
	require.NoError(t, err)

	// 创建配置
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建解析器
	p := parser.NewParser(cfg, log)

	// 解析项目
	endpoints, err := p.ParseProject(projectDir)
	require.NoError(t, err)

	// 创建 Swagger 构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 添加所有端点
	for _, endpoint := range endpoints {
		err := builder.AddEndpoint(endpoint)
		require.NoError(t, err)
	}

	// 构建文档
	doc := builder.Build()

	// 写入输出文件
	writer := output.NewWriter(outputDir)

	// 写入 YAML 格式
	err = writer.WriteSwagger(doc, "swagger", "yaml")
	require.NoError(t, err)

	// 验证文件
	swaggerFile := filepath.Join(outputDir, "swagger.yaml")
	assert.FileExists(t, swaggerFile)

	// 读取文件内容
	content, err := os.ReadFile(swaggerFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "openapi: 3.0.0")
	assert.Contains(t, string(content), "title: Test API")
}

// TestInitCommandIntegration_MultipleEndpoints 测试多个端点
func TestInitCommandIntegration_MultipleEndpoints(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	projectDir, err := os.MkdirTemp("", "init-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// 创建临时输出目录
	outputDir, err := os.MkdirTemp("", "init-output-*")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	// 创建多个 API 文件
	files := map[string]string{
		"user_api.go": `package api

// GetUsers 获取用户
// @Router /api/users [GET]
// @Summary 获取用户
// @Tags User
// @Success 200 {array} User
func GetUsers() {}

// CreateUser 创建用户
// @Router /api/users [POST]
// @Summary 创建用户
// @Tags User
// @Success 201 {object} User
func CreateUser() {}

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`,
		"post_api.go": `package api

// GetPosts 获取文章
// @Router /api/posts [GET]
// @Summary 获取文章
// @Tags Post
// @Success 200 {array} Post
func GetPosts() {}

// CreatePost 创建文章
// @Router /api/posts [POST]
// @Summary 创建文章
// @Tags Post
// @Success 201 {object} Post
func CreatePost() {}

type Post struct {
	ID    int    ` + "`json:\"id\"`" + `
	Title string ` + "`json:\"title\"`" + `
}
`,
	}

	for filename, content := range files {
		filePath := filepath.Join(projectDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// 创建配置
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建解析器
	p := parser.NewParser(cfg, log)

	// 解析项目
	endpoints, err := p.ParseProject(projectDir)
	require.NoError(t, err)
	assert.Equal(t, 4, len(endpoints), "应该找到 4 个端点")

	// 创建 Swagger 构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 添加所有端点
	for _, endpoint := range endpoints {
		err := builder.AddEndpoint(endpoint)
		require.NoError(t, err)
	}

	// 构建文档
	doc := builder.Build()
	assert.NotNil(t, doc)

	// 验证路径
	assert.Equal(t, 2, len(doc.Paths), "应该有 2 个路径")
	assert.Contains(t, doc.Paths, "/api/users")
	assert.Contains(t, doc.Paths, "/api/posts")

	// 验证方法
	usersPath := doc.Paths["/api/users"]
	assert.NotNil(t, usersPath.Get)
	assert.NotNil(t, usersPath.Post)

	postsPath := doc.Paths["/api/posts"]
	assert.NotNil(t, postsPath.Get)
	assert.NotNil(t, postsPath.Post)
}

// TestInitCommandIntegration_NestedDirectories 测试嵌套目录
func TestInitCommandIntegration_NestedDirectories(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	projectDir, err := os.MkdirTemp("", "init-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(projectDir)

	// 创建临时输出目录
	outputDir, err := os.MkdirTemp("", "init-output-*")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	// 创建嵌套目录
	userDir := filepath.Join(projectDir, "user")
	postDir := filepath.Join(projectDir, "post")
	err = os.MkdirAll(userDir, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(postDir, 0755)
	require.NoError(t, err)

	// 创建 API 文件
	userFile := filepath.Join(userDir, "api.go")
	userContent := `package user

// GetUser 获取用户
// @Router /api/users/{id} [GET]
// @Summary 获取用户
// @Tags User
// @Success 200 {object} User
func GetUser(id int) {}

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	err = os.WriteFile(userFile, []byte(userContent), 0644)
	require.NoError(t, err)

	postFile := filepath.Join(postDir, "api.go")
	postContent := `package post

// GetPost 获取文章
// @Router /api/posts/{id} [GET]
// @Summary 获取文章
// @Tags Post
// @Success 200 {object} Post
func GetPost(id int) {}

type Post struct {
	ID    int    ` + "`json:\"id\"`" + `
	Title string ` + "`json:\"title\"`" + `
}
`
	err = os.WriteFile(postFile, []byte(postContent), 0644)
	require.NoError(t, err)

	// 创建配置
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建解析器
	p := parser.NewParser(cfg, log)

	// 解析项目
	endpoints, err := p.ParseProject(projectDir)
	require.NoError(t, err)
	assert.Equal(t, 2, len(endpoints), "应该找到 2 个端点")

	// 创建 Swagger 构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 添加所有端点
	for _, endpoint := range endpoints {
		err := builder.AddEndpoint(endpoint)
		require.NoError(t, err)
	}

	// 构建文档
	doc := builder.Build()
	assert.NotNil(t, doc)

	// 验证路径
	assert.Contains(t, doc.Paths, "/api/users/{id}")
	assert.Contains(t, doc.Paths, "/api/posts/{id}")
}

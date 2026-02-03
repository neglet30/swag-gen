package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/neglet30/swag-gen/pkg/logger"
	"github.com/neglet30/swag-gen/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParserIntegration_ParseSimpleProject 测试解析简单项目
func TestParserIntegration_ParseSimpleProject(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	tmpDir, err := os.MkdirTemp("", "parser-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建示例 API 文件
	apiFile := filepath.Join(tmpDir, "api.go")
	apiContent := `package api

// GetUsers 获取所有用户
// @Router /api/users [GET]
// @Summary 获取所有用户
// @Description 获取系统中的所有用户
// @Tags User
// @Success 200 {array} User
// @Failure 500 {object} ErrorResponse
func GetUsers() {
}

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
			Name:    "Test API",
			Version: "1.0.0",
		},
	}

	// 创建解析器
	p := parser.NewParser(cfg, log)

	// 解析项目
	endpoints, err := p.ParseProject(tmpDir)
	require.NoError(t, err)

	// 验证结果
	assert.Greater(t, len(endpoints), 0, "应该找到至少一个端点")

	// 查找 GetUsers 端点
	var getUsersEndpoint *parser.Endpoint
	for _, ep := range endpoints {
		if ep.Path == "/api/users" {
			getUsersEndpoint = ep
			break
		}
	}

	assert.NotNil(t, getUsersEndpoint, "应该找到 GetUsers 端点")
	if getUsersEndpoint != nil {
		assert.Equal(t, "GET", getUsersEndpoint.Method)
		assert.Equal(t, "获取所有用户", getUsersEndpoint.Summary)
		assert.Equal(t, "获取系统中的所有用户", getUsersEndpoint.Description)
		assert.Contains(t, getUsersEndpoint.Tags, "User")
	}
}

// TestParserIntegration_ParseMultipleFiles 测试解析多个文件
func TestParserIntegration_ParseMultipleFiles(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	tmpDir, err := os.MkdirTemp("", "parser-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建多个 API 文件
	files := map[string]string{
		"user_api.go": `package api

// GetUsers 获取用户列表
// @Router /api/users [GET]
// @Summary 获取用户列表
// @Tags User
// @Success 200 {array} User
func GetUsers() {}

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`,
		"post_api.go": `package api

// GetPosts 获取文章列表
// @Router /api/posts [GET]
// @Summary 获取文章列表
// @Tags Post
// @Success 200 {array} Post
func GetPosts() {}

type Post struct {
	ID    int    ` + "`json:\"id\"`" + `
	Title string ` + "`json:\"title\"`" + `
}
`,
	}

	for filename, content := range files {
		filePath := filepath.Join(tmpDir, filename)
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
	endpoints, err := p.ParseProject(tmpDir)
	require.NoError(t, err)

	// 验证结果
	assert.Equal(t, 2, len(endpoints), "应该找到 2 个端点")

	// 验证端点
	paths := make(map[string]bool)
	for _, ep := range endpoints {
		paths[ep.Path] = true
	}

	assert.True(t, paths["/api/users"], "应该找到 /api/users 端点")
	assert.True(t, paths["/api/posts"], "应该找到 /api/posts 端点")
}

// TestParserIntegration_ParseWithNestedDirectories 测试解析嵌套目录
func TestParserIntegration_ParseWithNestedDirectories(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	tmpDir, err := os.MkdirTemp("", "parser-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建嵌套目录
	userDir := filepath.Join(tmpDir, "user")
	postDir := filepath.Join(tmpDir, "post")
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
	endpoints, err := p.ParseProject(tmpDir)
	require.NoError(t, err)

	// 验证结果
	assert.Equal(t, 2, len(endpoints), "应该找到 2 个端点")

	// 验证端点
	paths := make(map[string]bool)
	for _, ep := range endpoints {
		paths[ep.Path] = true
	}

	assert.True(t, paths["/api/users/{id}"], "应该找到 /api/users/{id} 端点")
	assert.True(t, paths["/api/posts/{id}"], "应该找到 /api/posts/{id} 端点")
}

// TestParserIntegration_ParseWithErrors 测试解析包含错误的项目
func TestParserIntegration_ParseWithErrors(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	tmpDir, err := os.MkdirTemp("", "parser-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建有效的 API 文件
	validFile := filepath.Join(tmpDir, "valid.go")
	validContent := `package api

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
	err = os.WriteFile(validFile, []byte(validContent), 0644)
	require.NoError(t, err)

	// 创建无效的 API 文件
	invalidFile := filepath.Join(tmpDir, "invalid.go")
	invalidContent := `package api

// 这是一个无效的 Go 文件
func GetPosts( {
	// 缺少参数
}
`
	err = os.WriteFile(invalidFile, []byte(invalidContent), 0644)
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

	// 解析项目 - 应该继续处理有效文件
	endpoints, err := p.ParseProject(tmpDir)
	require.NoError(t, err)

	// 验证结果 - 应该至少找到有效文件中的端点
	assert.Greater(t, len(endpoints), 0, "应该找到至少一个端点")
}

// TestParserIntegration_ParseEmptyProject 测试解析空项目
func TestParserIntegration_ParseEmptyProject(t *testing.T) {
	// 初始化日志
	err := logger.Init("info", "text")
	require.NoError(t, err)
	defer logger.Sync()

	log := logger.GetLogger()

	// 创建临时项目目录
	tmpDir, err := os.MkdirTemp("", "parser-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

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
	endpoints, err := p.ParseProject(tmpDir)
	require.NoError(t, err)

	// 验证结果
	assert.Equal(t, 0, len(endpoints), "空项目应该没有端点")
}

package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewParser(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}

	parser := NewParser(cfg, logger)

	assert.NotNil(t, parser)
	assert.Equal(t, cfg, parser.config)
	assert.Equal(t, logger, parser.logger)
}

func TestParserFindGoFiles(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}
	parser := NewParser(cfg, logger)

	// 创建临时目录
	tmpDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte("package main"), 0644)
	require.NoError(t, err)

	// 创建测试文件（应该被跳过）
	testTestFile := filepath.Join(tmpDir, "test_test.go")
	err = os.WriteFile(testTestFile, []byte("package main"), 0644)
	require.NoError(t, err)

	// 查找文件
	files, err := parser.findGoFiles(tmpDir)
	require.NoError(t, err)

	// 验证结果
	assert.Len(t, files, 1)
	assert.Equal(t, testFile, files[0])
}

func TestParserParseFile(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}
	parser := NewParser(cfg, logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

// @Router /api/users [GET]
// @Summary Get all users
// @Description Get all users from database
// @Tags User
func GetUsers() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	endpoints, err := parser.ParseFile(testFile)
	require.NoError(t, err)

	// 验证结果
	assert.Len(t, endpoints, 1)
	assert.Equal(t, "GET", endpoints[0].Method)
	assert.Equal(t, "/api/users", endpoints[0].Path)
	assert.Equal(t, "Get all users", endpoints[0].Summary)
}

func TestParserParseProject(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}
	parser := NewParser(cfg, logger)

	// 创建临时目录
	tmpDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

// @Router /api/users [GET]
// @Summary Get all users
func GetUsers() {
}

// @Router /api/users [POST]
// @Summary Create user
func CreateUser() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析项目
	endpoints, err := parser.ParseProject(tmpDir)
	require.NoError(t, err)

	// 验证结果
	assert.Len(t, endpoints, 2)
	assert.Equal(t, "GET", endpoints[0].Method)
	assert.Equal(t, "POST", endpoints[1].Method)
}

func TestParserParseProjectInvalidPath(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}
	parser := NewParser(cfg, logger)

	// 解析不存在的路径
	endpoints, err := parser.ParseProject("/nonexistent/path")

	assert.Error(t, err)
	assert.Nil(t, endpoints)
}

func TestParseComments(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}
	parser := NewParser(cfg, logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

// @Router /api/users [GET]
// @Summary Get all users
// @Description Get all users from database
// @Tags User
// @Deprecated
func GetUsers() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	endpoints, err := parser.ParseFile(testFile)
	require.NoError(t, err)

	// 验证结果
	assert.Len(t, endpoints, 1)
	endpoint := endpoints[0]
	assert.Equal(t, "GET", endpoint.Method)
	assert.Equal(t, "/api/users", endpoint.Path)
	assert.Equal(t, "Get all users", endpoint.Summary)
	assert.Equal(t, "Get all users from database", endpoint.Description)
	assert.True(t, endpoint.Deprecated)
}

func TestParseCommentsWithoutRouter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}
	parser := NewParser(cfg, logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

// @Summary Get all users
// @Description Get all users from database
func GetUsers() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	endpoints, err := parser.ParseFile(testFile)
	require.NoError(t, err)

	// 验证结果（没有 @Router 标签，应该返回空）
	assert.Len(t, endpoints, 0)
}

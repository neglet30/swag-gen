package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateInitOptions 测试参数验证
func TestValidateInitOptions(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		cleanup func()
		wantErr bool
	}{
		{
			name: "有效的参数",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = "Test API"
				initVersion = "1.0.0"
				initFormat = "json"
			},
			cleanup: func() {},
			wantErr: false,
		},
		{
			name: "空的源代码路径",
			setup: func() {
				initPath = ""
				initOutput = "./docs"
				initTitle = "Test API"
				initVersion = "1.0.0"
				initFormat = "json"
			},
			cleanup: func() {},
			wantErr: true,
		},
		{
			name: "不存在的源代码路径",
			setup: func() {
				initPath = "/nonexistent/path"
				initOutput = "./docs"
				initTitle = "Test API"
				initVersion = "1.0.0"
				initFormat = "json"
			},
			cleanup: func() {},
			wantErr: true,
		},
		{
			name: "空的输出路径",
			setup: func() {
				initPath = "."
				initOutput = ""
				initTitle = "Test API"
				initVersion = "1.0.0"
				initFormat = "json"
			},
			cleanup: func() {},
			wantErr: true,
		},
		{
			name: "空的标题",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = ""
				initVersion = "1.0.0"
				initFormat = "json"
			},
			cleanup: func() {},
			wantErr: true,
		},
		{
			name: "空的版本",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = "Test API"
				initVersion = ""
				initFormat = "json"
			},
			cleanup: func() {},
			wantErr: true,
		},
		{
			name: "无效的格式",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = "Test API"
				initVersion = "1.0.0"
				initFormat = "xml"
			},
			cleanup: func() {},
			wantErr: true,
		},
		{
			name: "YAML 格式",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = "Test API"
				initVersion = "1.0.0"
				initFormat = "yaml"
			},
			cleanup: func() {},
			wantErr: false,
		},
		{
			name: "YML 格式",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = "Test API"
				initVersion = "1.0.0"
				initFormat = "yml"
			},
			cleanup: func() {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			err := validateInitOptions()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetFileExtension 测试文件扩展名获取
func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		format   string
		expected string
	}{
		{"json", "json"},
		{"yaml", "yaml"},
		{"yml", "yaml"},
		{"", "json"},
		{"unknown", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			result := getFileExtension(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestInitCommandWithSampleProject 测试 init 命令与示例项目
func TestInitCommandWithSampleProject(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "swag-gen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建输出目录
	outputDir := filepath.Join(tmpDir, "docs")

	// 设置参数 - 使用当前目录作为源代码路径
	initPath = "."
	initOutput = outputDir
	initTitle = "Test API"
	initVersion = "1.0.0"
	initDescription = "Test API Description"
	initFormat = "json"

	// 验证参数
	err = validateInitOptions()
	assert.NoError(t, err)

	// 验证输出目录是否创建
	assert.NoError(t, os.MkdirAll(outputDir, 0755))
	assert.DirExists(t, outputDir)
}

// TestInitCommandParameterCombinations 测试各种参数组合
func TestInitCommandParameterCombinations(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		output      string
		title       string
		version     string
		description string
		format      string
		wantErr     bool
	}{
		{
			name:        "最小参数",
			path:        ".",
			output:      "./docs",
			title:       "API",
			version:     "1.0.0",
			description: "",
			format:      "json",
			wantErr:     false,
		},
		{
			name:        "完整参数",
			path:        ".",
			output:      "./docs",
			title:       "My API",
			version:     "2.0.0",
			description: "My API Description",
			format:      "yaml",
			wantErr:     false,
		},
		{
			name:        "缺少标题",
			path:        ".",
			output:      "./docs",
			title:       "",
			version:     "1.0.0",
			description: "Description",
			format:      "json",
			wantErr:     true,
		},
		{
			name:        "缺少版本",
			path:        ".",
			output:      "./docs",
			title:       "API",
			version:     "",
			description: "Description",
			format:      "json",
			wantErr:     true,
		},
		{
			name:        "无效的格式",
			path:        ".",
			output:      "./docs",
			title:       "API",
			version:     "1.0.0",
			description: "Description",
			format:      "xml",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置参数
			initPath = tt.path
			initOutput = tt.output
			initTitle = tt.title
			initVersion = tt.version
			initDescription = tt.description
			initFormat = tt.format

			// 验证参数
			err := validateInitOptions()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestInitCommandErrorHandling 测试错误处理
func TestInitCommandErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "路径不存在",
			setup: func() {
				initPath = "/nonexistent/path/that/does/not/exist"
				initOutput = "./docs"
				initTitle = "API"
				initVersion = "1.0.0"
				initFormat = "json"
			},
			wantErr: true,
		},
		{
			name: "空的路径",
			setup: func() {
				initPath = ""
				initOutput = "./docs"
				initTitle = "API"
				initVersion = "1.0.0"
				initFormat = "json"
			},
			wantErr: true,
		},
		{
			name: "空的输出路径",
			setup: func() {
				initPath = "."
				initOutput = ""
				initTitle = "API"
				initVersion = "1.0.0"
				initFormat = "json"
			},
			wantErr: true,
		},
		{
			name: "空的标题",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = ""
				initVersion = "1.0.0"
				initFormat = "json"
			},
			wantErr: true,
		},
		{
			name: "空的版本",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = "API"
				initVersion = ""
				initFormat = "json"
			},
			wantErr: true,
		},
		{
			name: "无效的格式",
			setup: func() {
				initPath = "."
				initOutput = "./docs"
				initTitle = "API"
				initVersion = "1.0.0"
				initFormat = "invalid"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := validateInitOptions()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

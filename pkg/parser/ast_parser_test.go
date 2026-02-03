package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewASTParser(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	astParser := NewASTParser(logger)

	assert.NotNil(t, astParser)
	assert.Equal(t, logger, astParser.logger)
}

func TestASTParserParseFile(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	astFile, err := astParser.ParseFile(testFile)
	require.NoError(t, err)

	// 验证结果
	assert.NotNil(t, astFile)
	assert.Equal(t, "main", astFile.Name.Name)
}

func TestASTParserParseFileInvalid(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

func main() {
	fmt.Println("Hello"
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件（应该失败）
	astFile, err := astParser.ParseFile(testFile)

	assert.Error(t, err)
	assert.Nil(t, astFile)
}

func TestASTParserParseDirectory(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时目录
	tmpDir := t.TempDir()

	// 创建测试文件
	testFile1 := filepath.Join(tmpDir, "test1.go")
	content1 := `package main

func Test1() {
}
`
	err := os.WriteFile(testFile1, []byte(content1), 0644)
	require.NoError(t, err)

	testFile2 := filepath.Join(tmpDir, "test2.go")
	content2 := `package main

func Test2() {
}
`
	err = os.WriteFile(testFile2, []byte(content2), 0644)
	require.NoError(t, err)

	// 解析目录
	files, err := astParser.ParseDirectory(tmpDir)
	require.NoError(t, err)

	// 验证结果
	assert.Len(t, files, 2)
}

func TestASTParserExtractFunctions(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

func Func1() {
}

func Func2() {
}

type MyType struct {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	astFile, err := astParser.ParseFile(testFile)
	require.NoError(t, err)

	// 提取函数
	functions := astParser.ExtractFunctions(astFile)

	// 验证结果
	assert.Len(t, functions, 2)
	assert.Equal(t, "Func1", functions[0].Name.Name)
	assert.Equal(t, "Func2", functions[1].Name.Name)
}

func TestASTParserExtractComments(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

// This is a comment
// Another comment
func TestFunc() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	astFile, err := astParser.ParseFile(testFile)
	require.NoError(t, err)

	// 提取函数
	functions := astParser.ExtractFunctions(astFile)
	require.Len(t, functions, 1)

	// 提取注释
	comments := astParser.ExtractComments(functions[0])

	// 验证结果
	assert.Len(t, comments, 2)
	assert.Contains(t, comments[0], "This is a comment")
	assert.Contains(t, comments[1], "Another comment")
}

func TestASTParserFindSwaggerTags(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	comments := []string{
		"// @Router /api/users [GET]",
		"// @Summary Get all users",
		"// @Tags User",
	}

	// 查找标签
	tags := astParser.FindSwaggerTags(comments)

	// 验证结果
	assert.Len(t, tags, 3)
	assert.Contains(t, tags, "@Router")
	assert.Contains(t, tags, "@Summary")
	assert.Contains(t, tags, "@Tags")
}

func TestASTParserValidateAST(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

func main() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	astFile, err := astParser.ParseFile(testFile)
	require.NoError(t, err)

	// 验证 AST
	err = astParser.ValidateAST(astFile)
	assert.NoError(t, err)
}

func TestASTParserGetPackageName(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package mypackage

func main() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	astFile, err := astParser.ParseFile(testFile)
	require.NoError(t, err)

	// 获取包名
	pkgName := astParser.GetPackageName(astFile)

	// 验证结果
	assert.Equal(t, "mypackage", pkgName)
}

func TestASTParserGetImports(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

import (
	"fmt"
	"os"
)

func main() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	astFile, err := astParser.ParseFile(testFile)
	require.NoError(t, err)

	// 获取导入
	imports := astParser.GetImports(astFile)

	// 验证结果
	assert.Len(t, imports, 2)
}

func TestASTParserGetStructs(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	astParser := NewASTParser(logger)

	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

type User struct {
	Name string
	Age  int
}

type Post struct {
	Title string
}

func main() {
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// 解析文件
	astFile, err := astParser.ParseFile(testFile)
	require.NoError(t, err)

	// 获取结构体
	structs := astParser.GetStructs(astFile)

	// 验证结果
	assert.Len(t, structs, 2)
	assert.Contains(t, structs, "User")
	assert.Contains(t, structs, "Post")
}

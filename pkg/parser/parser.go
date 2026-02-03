package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/neglet30/swag-gen/pkg/config"
	"go.uber.org/zap"
)

// Parser 代表代码解析器
type Parser struct {
	config *config.Config
	logger *zap.Logger
	mu     sync.Mutex
}

// NewParser 创建一个新的解析器
func NewParser(cfg *config.Config, logger *zap.Logger) *Parser {
	return &Parser{
		config: cfg,
		logger: logger,
	}
}

// ParseProject 解析整个项目
func (p *Parser) ParseProject(projectPath string) ([]*Endpoint, error) {
	p.logger.Info("开始解析项目", zap.String("path", projectPath))

	// 验证路径
	if _, err := os.Stat(projectPath); err != nil {
		p.logger.Error("项目路径不存在", zap.String("path", projectPath), zap.Error(err))
		return nil, fmt.Errorf("项目路径不存在: %w", err)
	}

	// 获取所有 Go 文件
	files, err := p.findGoFiles(projectPath)
	if err != nil {
		p.logger.Error("获取 Go 文件失败", zap.Error(err))
		return nil, fmt.Errorf("获取 Go 文件失败: %w", err)
	}

	p.logger.Info("找到 Go 文件", zap.Int("count", len(files)))

	// 并发解析文件
	endpoints := make([]*Endpoint, 0)
	var wg sync.WaitGroup
	resultChan := make(chan []*Endpoint, len(files))
	errChan := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			eps, err := p.ParseFile(filePath)
			if err != nil {
				p.logger.Warn("解析文件失败", zap.String("file", filePath), zap.Error(err))
				errChan <- err
				return
			}
			if len(eps) > 0 {
				resultChan <- eps
			}
		}(file)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	close(resultChan)
	close(errChan)

	// 收集结果
	for eps := range resultChan {
		endpoints = append(endpoints, eps...)
	}

	p.logger.Info("项目解析完成", zap.Int("endpoints", len(endpoints)))
	return endpoints, nil
}

// ParseFile 解析单个文件
func (p *Parser) ParseFile(filePath string) ([]*Endpoint, error) {
	p.logger.Debug("解析文件", zap.String("file", filePath))

	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 创建 FileSet
	fset := token.NewFileSet()

	// 解析 AST
	astFile, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		p.logger.Warn("解析 AST 失败", zap.String("file", filePath), zap.Error(err))
		return nil, fmt.Errorf("解析 AST 失败: %w", err)
	}

	// 提取 API 信息
	endpoints := p.extractEndpoints(astFile, filePath)
	p.logger.Debug("文件解析完成", zap.String("file", filePath), zap.Int("endpoints", len(endpoints)))

	return endpoints, nil
}

// findGoFiles 查找所有 Go 文件
func (p *Parser) findGoFiles(projectPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和非 Go 文件
		if info.IsDir() {
			// 跳过 vendor 和 test 目录
			if info.Name() == "vendor" || info.Name() == "test" {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) == ".go" && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// extractEndpoints 从 AST 中提取端点
func (p *Parser) extractEndpoints(file *ast.File, filePath string) []*Endpoint {
	var endpoints []*Endpoint

	// 遍历所有声明
	for _, decl := range file.Decls {
		// 查找函数声明
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// 检查函数是否有注释
		if funcDecl.Doc == nil {
			continue
		}

		// 解析注释
		endpoint := p.parseComments(funcDecl.Doc, filePath, int(funcDecl.Pos()))
		if endpoint != nil {
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

// parseComments 解析注释
func (p *Parser) parseComments(doc *ast.CommentGroup, filePath string, line int) *Endpoint {
	if doc == nil {
		return nil
	}

	endpoint := &Endpoint{
		File:       filePath,
		Line:       line,
		Tags:       make([]string, 0),
		Parameters: make([]Parameter, 0),
		Responses:  make(map[string]Response),
	}

	hasRouter := false

	// 遍历所有注释行
	for _, comment := range doc.List {
		text := comment.Text
		// 移除注释前缀 //
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimSpace(text)

		// 解析 @Router 标签
		if router := parseRouterTag(text); router != nil {
			endpoint.Method = router.Method
			endpoint.Path = router.Path
			hasRouter = true
		}

		// 解析 @Summary 标签
		if summary := parseTag(text, "@Summary"); summary != "" {
			endpoint.Summary = summary
		}

		// 解析 @Description 标签
		if desc := parseTag(text, "@Description"); desc != "" {
			endpoint.Description = desc
		}

		// 解析 @Tags 标签
		if tags := parseTag(text, "@Tags"); tags != "" {
			endpoint.Tags = append(endpoint.Tags, tags)
		}

		// 解析 @Param 标签
		if param := parseParamTag(text); param != nil {
			endpoint.Parameters = append(endpoint.Parameters, *param)
		}

		// 解析 @Success 标签
		if resp := parseResponseTag(text, "@Success"); resp != nil {
			endpoint.Responses[resp.StatusCode] = *resp
		}

		// 解析 @Failure 标签
		if resp := parseResponseTag(text, "@Failure"); resp != nil {
			endpoint.Responses[resp.StatusCode] = *resp
		}

		// 解析 @Deprecated 标签
		if strings.Contains(text, "@Deprecated") {
			endpoint.Deprecated = true
		}
	}

	// 只返回有 @Router 标签的端点
	if !hasRouter {
		return nil
	}

	return endpoint
}

// parseRouterTag 解析 @Router 标签
func parseRouterTag(text string) *struct {
	Method string
	Path   string
} {
	// 格式: @Router /api/users [GET]
	if !strings.Contains(text, "@Router") {
		return nil
	}

	// 使用正则表达式解析
	re := regexp.MustCompile(`@Router\s+(\S+)\s+\[(\w+)\]`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 3 {
		return nil
	}

	return &struct {
		Method string
		Path   string
	}{
		Method: matches[2],
		Path:   matches[1],
	}
}

// parseTag 解析简单标签
func parseTag(text, tag string) string {
	if !strings.Contains(text, tag) {
		return ""
	}

	// 提取标签后的内容
	idx := strings.Index(text, tag)
	if idx == -1 {
		return ""
	}

	content := text[idx+len(tag):]
	// 移除前导空格
	content = strings.TrimSpace(content)

	return content
}

// parseParamTag 解析 @Param 标签
func parseParamTag(text string) *Parameter {
	if !strings.Contains(text, "@Param") {
		return nil
	}

	// 简单的解析逻辑
	// 格式: @Param page query int false "Page number"
	return nil
}

// parseResponseTag 解析响应标签
func parseResponseTag(text, tag string) *Response {
	if !strings.Contains(text, tag) {
		return nil
	}

	// 简单的解析逻辑
	// 格式: @Success 200 {object} User
	return nil
}

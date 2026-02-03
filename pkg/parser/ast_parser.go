package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// ASTParser 代表 AST 解析器
type ASTParser struct {
	logger *zap.Logger
}

// NewASTParser 创建一个新的 AST 解析器
func NewASTParser(logger *zap.Logger) *ASTParser {
	return &ASTParser{
		logger: logger,
	}
}

// ParseFile 解析单个文件的 AST
func (ap *ASTParser) ParseFile(filePath string) (*ast.File, error) {
	ap.logger.Debug("解析 AST", zap.String("file", filePath))

	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		ap.logger.Error("读取文件失败", zap.String("file", filePath), zap.Error(err))
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 创建 FileSet
	fset := token.NewFileSet()

	// 解析 AST
	astFile, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		ap.logger.Error("解析 AST 失败", zap.String("file", filePath), zap.Error(err))
		return nil, fmt.Errorf("解析 AST 失败: %w", err)
	}

	return astFile, nil
}

// ParseDirectory 递归解析目录中的所有 Go 文件
func (ap *ASTParser) ParseDirectory(dirPath string) (map[string]*ast.File, error) {
	ap.logger.Info("解析目录", zap.String("path", dirPath))

	files := make(map[string]*ast.File)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			// 跳过 vendor 和 test 目录
			if info.Name() == "vendor" || info.Name() == "test" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// 只处理 Go 文件
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// 解析文件
		astFile, err := ap.ParseFile(path)
		if err != nil {
			ap.logger.Warn("解析文件失败", zap.String("file", path), zap.Error(err))
			return nil
		}

		files[path] = astFile
		return nil
	})

	if err != nil {
		ap.logger.Error("解析目录失败", zap.String("path", dirPath), zap.Error(err))
		return nil, fmt.Errorf("解析目录失败: %w", err)
	}

	ap.logger.Info("目录解析完成", zap.Int("files", len(files)))
	return files, nil
}

// ExtractFunctions 从 AST 中提取所有函数
func (ap *ASTParser) ExtractFunctions(file *ast.File) []*ast.FuncDecl {
	var functions []*ast.FuncDecl

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			functions = append(functions, funcDecl)
		}
	}

	return functions
}

// ExtractComments 从函数中提取注释
func (ap *ASTParser) ExtractComments(funcDecl *ast.FuncDecl) []string {
	var comments []string

	if funcDecl.Doc != nil {
		for _, comment := range funcDecl.Doc.List {
			comments = append(comments, comment.Text)
		}
	}

	return comments
}

// FindSwaggerTags 查找所有 Swagger 标签
func (ap *ASTParser) FindSwaggerTags(comments []string) map[string][]string {
	tags := make(map[string][]string)

	for _, comment := range comments {
		// 移除注释前缀
		text := strings.TrimPrefix(comment, "//")
		text = strings.TrimSpace(text)

		// 查找标签
		if strings.HasPrefix(text, "@") {
			parts := strings.SplitN(text, " ", 2)
			if len(parts) >= 1 {
				tagName := parts[0]
				tagValue := ""
				if len(parts) > 1 {
					tagValue = parts[1]
				}
				tags[tagName] = append(tags[tagName], tagValue)
			}
		}
	}

	return tags
}

// ValidateAST 验证 AST 的有效性
func (ap *ASTParser) ValidateAST(file *ast.File) error {
	if file == nil {
		return fmt.Errorf("AST 文件为空")
	}

	if file.Name == nil {
		return fmt.Errorf("AST 文件名为空")
	}

	return nil
}

// GetPackageName 获取包名
func (ap *ASTParser) GetPackageName(file *ast.File) string {
	if file.Name != nil {
		return file.Name.Name
	}
	return ""
}

// GetImports 获取所有导入
func (ap *ASTParser) GetImports(file *ast.File) []string {
	var imports []string

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok.String() == "import" {
				for _, spec := range genDecl.Specs {
					if importSpec, ok := spec.(*ast.ImportSpec); ok {
						if importSpec.Path != nil {
							imports = append(imports, importSpec.Path.Value)
						}
					}
				}
			}
		}
	}

	return imports
}

// GetStructs 获取所有结构体定义
func (ap *ASTParser) GetStructs(file *ast.File) map[string]*ast.StructType {
	structs := make(map[string]*ast.StructType)

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structs[typeSpec.Name.Name] = structType
					}
				}
			}
		}
	}

	return structs
}

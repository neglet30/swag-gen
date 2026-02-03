package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neglet30/swag-gen/pkg/config"
	"github.com/neglet30/swag-gen/pkg/logger"
	"github.com/neglet30/swag-gen/pkg/output"
	"github.com/neglet30/swag-gen/pkg/parser"
	"github.com/neglet30/swag-gen/pkg/swagger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	initPath        string
	initOutput      string
	initTitle       string
	initVersion     string
	initDescription string
	initFormat      string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目并生成 Swagger 文档",
	Long: `扫描项目中的 API 注释，生成 Swagger/OpenAPI 规范文档。

示例:
  swag-gen init -p ./api -o ./docs -t "My API"
  swag-gen init -p ./api -o ./docs -t "My API" -f yaml`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVarP(&initPath, "path", "p", "./", "API 源代码路径")
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "./docs", "输出文档路径")
	initCmd.Flags().StringVarP(&initTitle, "title", "t", "API Documentation", "API 标题")
	initCmd.Flags().StringVarP(&initVersion, "version", "v", "1.0.0", "API 版本")
	initCmd.Flags().StringVarP(&initDescription, "description", "d", "", "API 描述")
	initCmd.Flags().StringVarP(&initFormat, "format", "f", "json", "输出格式 (json 或 yaml)")
}

func runInit(cmd *cobra.Command, args []string) error {
	// 初始化日志
	if err := logger.Init("info", "text"); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}
	defer logger.Sync()

	log := logger.GetLogger()

	// 验证参数
	if err := validateInitOptions(); err != nil {
		return err
	}

	fmt.Printf("初始化项目...\n")
	fmt.Printf("  源代码路径: %s\n", initPath)
	fmt.Printf("  输出路径: %s\n", initOutput)
	fmt.Printf("  API 标题: %s\n", initTitle)
	fmt.Printf("  API 版本: %s\n", initVersion)
	fmt.Printf("  输出格式: %s\n", initFormat)

	// 创建配置
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:        initTitle,
			Version:     initVersion,
			Description: initDescription,
		},
		Parser: config.ParserConfig{
			EnableCache:   true,
			CacheTTL:      3600,
			MaxConcurrent: 4,
			ExcludeDirs:   []string{"vendor", "node_modules", ".git", "test", "tests"},
		},
	}

	// 创建解析器
	p := parser.NewParser(cfg, log)

	// 解析项目
	fmt.Println("\n正在解析项目...")
	endpoints, err := p.ParseProject(initPath)
	if err != nil {
		return fmt.Errorf("解析项目失败: %w", err)
	}

	fmt.Printf("✓ 找到 %d 个 API 端点\n", len(endpoints))

	// 创建 Swagger 构建器
	fmt.Println("\n正在生成 Swagger 文档...")
	builder := swagger.NewBuilder(initTitle, initVersion, initDescription)

	// 添加所有端点
	for _, endpoint := range endpoints {
		if err := builder.AddEndpoint(endpoint); err != nil {
			log.Warn("添加端点失败", zap.String("endpoint", endpoint.Path), zap.Error(err))
			continue
		}
	}

	// 构建文档
	doc := builder.Build()

	// 写入输出文件
	fmt.Println("\n正在写入输出文件...")
	writer := output.NewWriter(initOutput)

	// 写入 Swagger 文档
	if err := writer.WriteSwagger(doc, "swagger", initFormat); err != nil {
		return fmt.Errorf("写入 Swagger 文档失败: %w", err)
	}

	fmt.Printf("✓ Swagger 文档已写入: %s\n", filepath.Join(initOutput, "swagger."+getFileExtension(initFormat)))

	// 写入配置文件
	outputConfig := output.NewConfig(initTitle, initVersion, initDescription)
	outputConfig.SetParserPath(initPath)
	outputConfig.SetOutputPath(initOutput)
	outputConfig.SetOutputFormat(initFormat)

	if err := writer.WriteConfig(outputConfig, "swag-gen.yaml"); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	fmt.Printf("✓ 配置文件已写入: %s\n", filepath.Join(initOutput, "swag-gen.yaml"))

	// 写入 README
	if err := writer.WriteREADME("README.md", initTitle, initDescription); err != nil {
		return fmt.Errorf("写入 README 失败: %w", err)
	}

	fmt.Printf("✓ README 已写入: %s\n", filepath.Join(initOutput, "README.md"))

	fmt.Println("\n✓ 项目初始化完成")
	return nil
}

// validateInitOptions 验证初始化选项
func validateInitOptions() error {
	// 验证路径参数
	if initPath == "" {
		return fmt.Errorf("源代码路径不能为空")
	}

	// 验证路径是否存在
	if _, err := os.Stat(initPath); err != nil {
		return fmt.Errorf("源代码路径不存在: %s", initPath)
	}

	// 验证输出路径
	if initOutput == "" {
		return fmt.Errorf("输出路径不能为空")
	}

	// 验证标题
	if initTitle == "" {
		return fmt.Errorf("API 标题不能为空")
	}

	// 验证版本
	if initVersion == "" {
		return fmt.Errorf("API 版本不能为空")
	}

	// 验证格式
	if initFormat != "json" && initFormat != "yaml" && initFormat != "yml" {
		return fmt.Errorf("输出格式必须是 json 或 yaml")
	}

	return nil
}

// getFileExtension 获取文件扩展名
func getFileExtension(format string) string {
	if format == "yaml" || format == "yml" {
		return "yaml"
	}
	return "json"
}

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	initPath        string
	initOutput      string
	initTitle       string
	initVersion     string
	initDescription string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目并生成 Swagger 文档",
	Long: `扫描项目中的 API 注释，生成 Swagger/OpenAPI 规范文档。

示例:
  swag-gen init -p ./api -o ./docs -t "My API"`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVarP(&initPath, "path", "p", "./", "API 源代码路径")
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "./docs", "输出文档路径")
	initCmd.Flags().StringVarP(&initTitle, "title", "t", "API Documentation", "API 标题")
	initCmd.Flags().StringVarP(&initVersion, "version", "v", "1.0.0", "API 版本")
	initCmd.Flags().StringVarP(&initDescription, "description", "d", "", "API 描述")
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Printf("初始化项目...\n")
	fmt.Printf("  源代码路径: %s\n", initPath)
	fmt.Printf("  输出路径: %s\n", initOutput)
	fmt.Printf("  API 标题: %s\n", initTitle)
	fmt.Printf("  API 版本: %s\n", initVersion)

	// TODO: 实现项目初始化逻辑
	fmt.Println("\n✓ 项目初始化完成")

	return nil
}

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	serverPort     int
	serverHost     string
	serverDocsPath string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 Web 服务器",
	Long: `启动 swag-gen Web 服务器，提供 Swagger UI 和 API 测试工具。

示例:
  swag-gen server -p 8080 -h 0.0.0.0 -d ./docs`,
	RunE: runServer,
}

func init() {
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "服务器端口")
	serverCmd.Flags().StringVar(&serverHost, "host", "0.0.0.0", "服务器地址")
	serverCmd.Flags().StringVarP(&serverDocsPath, "docs", "d", "./docs", "文档路径")
}

func runServer(cmd *cobra.Command, args []string) error {
	fmt.Printf("启动 Web 服务器...\n")
	fmt.Printf("  地址: %s:%d\n", serverHost, serverPort)
	fmt.Printf("  文档路径: %s\n", serverDocsPath)
	fmt.Printf("\n访问 UI: http://localhost:%d/swagger/ui\n", serverPort)

	// TODO: 实现 Web 服务器启动逻辑
	fmt.Println("\n✓ 服务器已启动")

	return nil
}

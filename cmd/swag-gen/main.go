package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	commit  = "unknown"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "swag-gen",
	Short: "API 文档生成工具",
	Long: `swag-gen 是一个开源的 API 文档生成工具库，功能与 swag 相同，
但提供更强大的 Web UI 界面。可以轻松集成到任何 Go 项目中。`,
	Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(serverCmd)
}

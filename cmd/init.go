package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"proxyCrawler/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置文件",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig(args[0])
		fmt.Printf("配置文件已存在,请使用[%s start all -c %s]启动服务", os.Args[0], configFlag)
		return
	},
}

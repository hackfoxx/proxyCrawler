package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"proxyCrawler/config"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "测试配置文件有效性",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if config.Check(args[0]) {
			fmt.Println("配置文件有效, Web.Authorization的值为[", config.Cfg.Web.Authorization, "]")
		} else {
			fmt.Println("请检查配置文件！！！")
		}
	},
}

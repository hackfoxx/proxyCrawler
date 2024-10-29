package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"proxyCrawler/config"
	"proxyCrawler/database"
	"proxyCrawler/jobs"
	"proxyCrawler/module/crawler"
	"proxyCrawler/utils/logger"
	"proxyCrawler/web"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动对应的服务,可选[web|crawler|all|]",
	Long:  fmt.Sprintf("Start Server, exp %s start web\n%s start", os.Args[0], os.Args[0]),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start(args[0])
	},
}

func start(options string) {
	// 初始化配置文件
	config.Init(configFlag)
	// 初始化logger
	logger.Init()
	// 初始化数据库
	database.Init()
	switch options {
	case "web", "w":
		{
			log.Println("webserver starting at " + "http://localhost" + config.Cfg.Web.Addr + config.Cfg.Web.BasePath)
			// 启动web服务
			web.NewGinServer().Start()
			select {}
		}
	case "crawler", "c":
		{
			// 初始化爬取器
			crawler.Init()
			// 启动调度器
			jobs.CrawlersJob()
			break
		}
	case "validator", "v":
		{
			crawler.Init()
			jobs.ValidatorJob()
			break
		}
	case "all", "a":
		{
			crawler.Init()
			log.Println("webserver starting at " + "http://localhost" + config.Cfg.Web.Addr + config.Cfg.Web.BasePath)
			go web.NewGinServer().Start()
			<-jobs.GetScheduler().Start()
			break
		}
	default:
		{
			fmt.Printf("命令错误,请使用: \n1. %s start web\n2. %s start crawler\n3. %s start validator\n4. %s start all\n或首字母简写: %s start w", os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
			os.Exit(0)
		}
	}
}

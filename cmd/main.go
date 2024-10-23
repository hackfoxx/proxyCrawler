package main

import (
	"fmt"
	"proxyCrawler/config"
	"proxyCrawler/database"
	"proxyCrawler/jobs"
	"proxyCrawler/module/crawler"
	"proxyCrawler/utils/logger"
	"proxyCrawler/web"
)

func _init() {
	fmt.Println("initializing logger....")
	if logger.Init() {
		fmt.Println("initialized logger...")
	}
	fmt.Println("initializing database....")
	if database.Init() != nil {
		fmt.Println("initialized database...")
	}
	fmt.Println("initializing crawlers....")
	if crawler.Init() {
		fmt.Println("initialized crawlers...")
	}
}
func start() {
	fmt.Println("webserver starting at " + config.GetConfig().Web.Addr)
	go web.NewGinServer().Start()
	<-jobs.GetScheduler().Start()
}
func main() {
	_init()
	//test()
	start()
}

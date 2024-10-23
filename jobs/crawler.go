package jobs

import (
	"fmt"
	"proxyCrawler/database"
	"proxyCrawler/model"
	"proxyCrawler/module/crawler"
	"proxyCrawler/utils"
	"proxyCrawler/utils/logger"
	"strconv"
	"time"
)

func CrawlersJob() bool {
	logger.GetLogger().Info("Crawler service started ...")
	// 从数据库获取所有的爬取器
	dbCrawler := database.GetCrawlers(nil)

	logger.GetLogger().Info("共发现 " + strconv.Itoa(len(dbCrawler)) + " 个爬取器")
	// 运行的代理数量
	proxiesCount := model.DBResult{}
	crawlerCount := model.DBResult{
		Sum: int64(len(dbCrawler)),
	}
	for _, Crawler := range dbCrawler {
		// 如果爬取器启用 && （最后一次爬取时间+一小时）小于 现在
		if Crawler.IsRunning {
			logger.GetLogger().Info("爬取器 " + Crawler.Type + " 正在运行，本次已跳过！")
		} else if !Crawler.Enable {
			logger.GetLogger().Info("爬取器 " + Crawler.Type + " 未启用，本次已跳过！")
		} else if Crawler.ToFetchDate.After(time.Now()) {
			logger.GetLogger().Info("爬取器 " + Crawler.Type + " 或未到执行时间，本次已跳过！")
		} else {
			tmp := Crawler
			tmp.IsRunning = true
			database.SetCrawlers([]model.Crawler{tmp}, true)
			logger.GetLogger().Info("爬取器 " + Crawler.Type + " 正在运行...")
			// 从网络或AddedData获取proxies
			craw := crawler.GetCrawler(Crawler.Type)
			if craw == nil {
				logger.GetLogger().Error("Error: 不存在" + Crawler.Type + "扫描器,请检查输入")
				return false
			}
			proxies := craw.Run()
			// 运行 crawler 爬取proxies
			proxiesResult := database.SetProxies(proxies, false)
			proxiesCount = utils.AddDBResult(proxiesCount, proxiesResult)
			// 最后抓取的数量 = crawler返回的代理数量
			Crawler.LastProxiesCnt = int(proxiesResult.Sum)
			// 代理总数+=数据库新增的代理数量
			Crawler.SumProxiesCnt += int(proxiesResult.Added)
			// 最后抓取的时间=现在
			Crawler.LastFetchDate = time.Now()
			// 下次计划时间
			Crawler.ToFetchDate = time.Now().Add(crawler.GetInterval(Crawler.Type))
			Crawler.IsRunning = false
			// 更新爬取器
			crawlerCount = utils.AddDBResult(crawlerCount, database.SetCrawlers([]model.Crawler{Crawler}, true))
			logger.GetLogger().Info(fmt.Sprintf("爬取器 %s 运行结束，本次共发现 %d，更新 %d,新增 %d 个代理\n", Crawler.Type, proxiesResult.Sum, proxiesResult.Updated, proxiesResult.Added))
		}
	}
	logger.GetLogger().Info(fmt.Sprintf("本次共发现%d个爬取器，运行%d个。共发现 %d，更新 %d,新增 %d 个代理", len(dbCrawler), crawlerCount.Updated, proxiesCount.Sum, proxiesCount.Updated, proxiesCount.Added))
	return true
}

func CrawlerJob(cType string) string {
	logger.GetLogger().Info("正在尝试手动运行爬取器 " + cType)
	Crawler := database.GetCrawler(cType)
	if Crawler.Type != cType {
		logger.GetLogger().Info("爬取器" + cType + "不存在")
		return "爬取器" + cType + "不存在"
	}
	// 如果爬取器启用 && （最后一次爬取时间+一小时）小于 现在
	if Crawler.IsRunning {
		logger.GetLogger().Info("爬取器 " + Crawler.Type + " 正在运行，本次已跳过！")
		return "爬取器 " + Crawler.Type + " 正在运行，本次已跳过！"
	} else if !Crawler.Enable {
		logger.GetLogger().Info("爬取器 " + Crawler.Type + " 未启用，本次已跳过！")
		return "爬取器 " + Crawler.Type + " 未启用，本次已跳过！"
	} else if Crawler.ToFetchDate.After(time.Now()) {
		logger.GetLogger().Info("爬取器 " + Crawler.Type + " 或未到执行时间，本次已跳过！")
		return "爬取器 " + Crawler.Type + " 或未到执行时间，本次已跳过！"
	} else {
		tmp := Crawler
		tmp.IsRunning = true
		database.SetCrawlers([]model.Crawler{tmp}, true)
		logger.GetLogger().Info("爬取器 " + Crawler.Type + " 正在运行...")
		// 从网络或AddedData获取proxies
		craw := crawler.GetCrawler(Crawler.Type)
		if craw == nil {
			logger.GetLogger().Error("Error: 不存在" + Crawler.Type + "扫描器,请检查输入")
			return "Error: 不存在" + Crawler.Type + "扫描器,请检查输入"
		}
		proxies := craw.Run()
		// 运行 crawler 爬取proxies
		proxiesResult := database.SetProxies(proxies, false)
		// 最后抓取的数量 = crawler返回的代理数量
		Crawler.LastProxiesCnt = int(proxiesResult.Sum)
		// 代理总数+=数据库新增的代理数量
		Crawler.SumProxiesCnt += int(proxiesResult.Added)
		// 最后抓取的时间=现在
		Crawler.LastFetchDate = time.Now()
		// 下次计划时间
		Crawler.ToFetchDate = time.Now().Add(crawler.GetInterval(Crawler.Type))
		Crawler.IsRunning = false
		// 更新爬取器
		database.SetCrawlers([]model.Crawler{Crawler}, true)
		logger.GetLogger().Info(fmt.Sprintf("爬取器 %s 运行结束，本次共发现 %d，更新 %d,新增 %d 个代理\n", Crawler.Type, proxiesResult.Sum, proxiesResult.Updated, proxiesResult.Added))
		return fmt.Sprintf("爬取器 %s 运行结束，本次共发现 %d，更新 %d,新增 %d 个代理\n", Crawler.Type, proxiesResult.Sum, proxiesResult.Updated, proxiesResult.Added)
	}
}

package crawler

import (
	"proxyCrawler/database"
	"proxyCrawler/model"
	"time"
)

// Crawler 爬取器接口
type Crawler interface {
	Run() []model.Proxy
}

// crawlerMap 用于注册Crawler
var crawlerMap = map[string]Crawler{
	"ip3366": ip3366Crawler{},
	"local":  localCrawler{},
	"xui":    xuiCrawler{},
	"89ip":   c89ipCrawler{},
	"kx":     KxCrawler{},
}

var intervalMap = map[string]time.Duration{
	"ip3366": time.Hour * 4,
	"local":  0,
	"xui":    0,
	"89ip":   time.Hour * 4,
	"kx":     time.Hour * 4,
}

func GetInterval(ctype string) time.Duration {

	return intervalMap[ctype]
}
func Init() bool {
	for typeName := range crawlerMap {
		crawler := database.GetCrawler(typeName)
		tmp := crawler
		if crawler.Type != typeName {
			database.SetCrawlers([]model.Crawler{model.Crawler{Type: typeName, Enable: true, IsRunning: false}}, false)
			continue
		} else if crawler.IsRunning {
			tmp.IsRunning = false
			database.SetCrawlers([]model.Crawler{tmp}, true)
		}
	}
	return true
}

// GetCrawlerTypes 获取全部CrawlerType
func GetCrawlerTypes() []string {
	var keys []string
	for key := range crawlerMap {
		keys = append(keys, key)
	}
	return keys
}

// GetCrawler 通过注册的名称获取对应的类型
func GetCrawler(typeName string) Crawler {
	return crawlerMap[typeName]
}

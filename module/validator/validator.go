package validator

import (
	"fmt"
	"proxyCrawler/database"
	"proxyCrawler/model"
	"proxyCrawler/module/crawler"
	"proxyCrawler/utils"
	"proxyCrawler/utils/logger"
	"strings"
	"sync"
	"time"
)

func getToValidateDate(proxy model.Proxy) time.Time {
	return proxy.ValidateDate.Add(crawler.GetInterval(proxy.CType))
}

// ValidProxy 验证器
func ValidProxy(proxy model.Proxy) model.Proxy {
	latency := checkValid(proxy)
	if latency != 0 {
		proxy.Validated = true
		proxy.ValidateFailedCnt = 0
		proxy.Latency = int(latency)
	} else {
		proxy.Validated = false
		proxy.ValidateFailedCnt++
		proxy.Latency = 99999
	}
	if proxy.Link == "" && proxy.Protocol != "vmess" {
		proxy.Link = fmt.Sprintf("%s://%s:%s@%s:%d", proxy.Protocol, proxy.User, proxy.Pass, proxy.Host, proxy.Port)
		proxy.Link = strings.ReplaceAll(proxy.Link, "://:@", "://")
	}
	proxy.ValidateDate = time.Now()
	proxy.ToValidateDate = getToValidateDate(proxy)
	proxy.Country = crawler.GetIPCountry(proxy.Host)
	return proxy
}

func ValidProxies(proxies []model.Proxy) model.DBResult {
	var result model.DBResult
	var mu sync.Mutex // 互斥锁，用于保护 result
	wg := sync.WaitGroup{}
	ch := make(chan model.DBResult, len(proxies)) // 通道缓冲区大小与代理数量一致

	concurrencyLimit := 100 // 设置并发限制数
	sem := make(chan struct{}, concurrencyLimit)
	progressChan := make(chan int)
	for _, proxy := range proxies {
		wg.Add(1)
		go func(proxy model.Proxy, ch chan model.DBResult, wg *sync.WaitGroup, sem chan struct{}, progress chan int) {
			defer wg.Done()
			defer func() { progress <- 1 }()
			sem <- struct{}{}        // 获取一个信号量
			defer func() { <-sem }() // 释放信号量
			count := 3
			var dbResult model.DBResult
			for dbResult = database.SetProxy(ValidProxy(proxy), true); dbResult.Error == 1 && count > 0; count-- {
				logger.GetLogger().Info("正在重试...")
				dbResult = database.SetProxy(ValidProxy(proxy), true)
			}
			mu.Lock()
			result = utils.AddDBResult(result, dbResult)
			mu.Unlock()
		}(proxy, ch, &wg, sem, progressChan)
	}
	go func() {
		wg.Wait()
		close(progressChan)
	}()
	count := 0
	for i := range progressChan {
		count = count + i
		fmt.Printf("验证进度： %d / %d\n", count, len(proxies))
	}
	return result
}

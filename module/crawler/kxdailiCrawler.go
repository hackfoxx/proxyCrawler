package crawler

import (
	"fmt"
	"io"
	"net/http"
	"proxyCrawler/model"
	"proxyCrawler/utils"
	"proxyCrawler/utils/logger"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type KxCrawler struct{}

func (receiver KxCrawler) Run() []model.Proxy {
	var result []model.Proxy
	wg := sync.WaitGroup{}
	ch := make(chan model.Proxy, 200)
	for p1 := 1; p1 <= 2; p1++ {
		for p2 := 1; p2 <= 10; p2++ {
			wg.Add(1)
			go GetPage(p1, p2, &wg, ch)
		}
	}
	wg.Wait()
	close(ch)
	for proxy := range ch {
		result = append(result, proxy)
	}
	return result
}
func GetPage(p1, p2 int, wg *sync.WaitGroup, ch chan model.Proxy) []model.Proxy {
	defer wg.Done()
	request, err := http.NewRequest("GET", fmt.Sprintf("http://www.kxdaili.com/dailiip/%d/%d.html", p1, p2), nil)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	client, _ := utils.GetClient(false)
	response, err := client.Do(request)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil
	}
	defer response.Body.Close()
	websitePage := string(body)
	re := regexp.MustCompile(`<td>([\d.]+)</td>\s*<td>(\d+)</td>\s*<td>.*?</td>\s*<td>((\w|,)+)</td>`)
	matches := re.FindAllStringSubmatch(websitePage, -1)
	var result []model.Proxy
	for _, match := range matches {
		ip := match[1]
		port, _ := strconv.Atoi(match[2])
		protocol := "http"
		if strings.Contains(match[3], "HTTPS") {
			protocol = "https"
		}
		ch <- model.Proxy{
			CType:             "kx",
			Protocol:          protocol,
			Host:              ip,
			Port:              port,
			Validated:         false,
			Latency:           0,
			ValidateDate:      time.Time{},
			ToValidateDate:    time.Time{},
			ValidateFailedCnt: 0,
			User:              "",
			Pass:              "",
			Link:              "",
			Country:           "",
		}
	}
	return result
}

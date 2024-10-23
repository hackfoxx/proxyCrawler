package crawler

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"proxyCrawler/model"
	"proxyCrawler/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ip3366Crawler struct {
}

func (f ip3366Crawler) Run() []model.Proxy {
	var proxies []model.Proxy
	maxStype := 2
	maxPage := 7
	wg := sync.WaitGroup{}
	proxyChan := make(chan model.Proxy, 400)
	for stype := 1; stype <= maxStype; stype++ {
		for page := 1; page <= maxPage; page++ {
			wg.Add(1)
			go workflow(stype, page, &wg, proxyChan)
		}
	}
	wg.Wait()
	close(proxyChan)
	for proxy := range proxyChan {
		proxies = append(proxies, proxy)
	}
	return proxies
}
func workflow(stype, page int, wg *sync.WaitGroup, proxyChan chan model.Proxy) {
	defer wg.Done()
	client, _ := utils.GetClient(true)
	resp, err := client.Get(fmt.Sprintf("http://www.ip3366.net/free/?stype=%d&page=%d", stype, page))
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println(err.Error())
		return
	}
	websitePage := string(respBody)
	doc, err := htmlquery.Parse(strings.NewReader(websitePage))
	trs, err := htmlquery.QueryAll(doc, "/html/body/div[2]/div/div[2]/table/tbody/tr")
	if err != nil {
		//fmt.Println(err)
		return
	}
	for _, tr := range trs {
		tds, err := htmlquery.QueryAll(tr, "/td")
		if err != nil {
			return
		}
		port, _ := strconv.Atoi(htmlquery.InnerText(tds[1]))
		var proxy = model.Proxy{
			CType:             "ip3366",
			Protocol:          strings.ToLower(htmlquery.InnerText(tds[3])),
			Host:              htmlquery.InnerText(tds[0]),
			Port:              port,
			Validated:         false,
			Latency:           0,
			ValidateDate:      time.Time{},
			ToValidateDate:    time.Time{},
			ValidateFailedCnt: 0,
			User:              "",
			Pass:              "",
			Link:              "",
			Country:           "CN",
		}
		proxyChan <- proxy
	}
}

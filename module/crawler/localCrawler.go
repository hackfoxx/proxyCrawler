package crawler

import (
	"log"
	"net/url"
	"proxyCrawler/database"
	"proxyCrawler/model"
	"strconv"
	"time"
)

type localCrawler struct {
}

// todo
func (f localCrawler) Run() []model.Proxy {
	var result []model.Proxy
	proxies := database.GetAddedData(map[string]interface{}{"c_type": "local", "validated": false})
	for _, proxy := range proxies {
		database.UpdateAddedDataValidated(model.AddedData{CType: "local", Data: proxy.Data})
		// 解析protocol://username:password@host:port
		parsedURL, err := url.Parse(proxy.Data)
		if err != nil {
			log.Fatal(err)
		}
		protocol := parsedURL.Scheme
		username := ""
		password := ""
		if parsedURL.User != nil {
			username = parsedURL.User.Username()
			pw, isSet := parsedURL.User.Password()
			if isSet {
				password = pw
			}
		}
		host := parsedURL.Hostname()
		port, _ := strconv.Atoi(parsedURL.Port())
		node := model.Proxy{
			CType:             "local",
			Protocol:          protocol,
			Host:              host,
			Port:              port,
			Validated:         false,
			Latency:           0,
			ValidateDate:      time.Time{},
			ToValidateDate:    time.Time{},
			ValidateFailedCnt: 0,
			User:              username,
			Pass:              password,
			Link:              "",
			Country:           "",
		}
		result = append(result, node)
	}
	return result
}

package crawler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"proxyCrawler/database"
	"proxyCrawler/model"
	"proxyCrawler/utils"
	"proxyCrawler/utils/logger"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type xuiCrawler struct{}

func (f xuiCrawler) Run() []model.Proxy {
	urls := database.GetAddedData(map[string]interface{}{"c_type": "xui", "validated": false})
	var proxies []model.Proxy
	var wg sync.WaitGroup
	ch := make(chan model.Proxy, len(urls)*10)
	concurrencyLimit := 100 // 设置并发限制数
	count := len(urls)
	if count > 5000 {
		count = 5000
	}
	logger.GetLogger().Info(fmt.Sprintf("xui: 发现 %d 个代理待爬取,本次计划爬取 %d 个", len(urls), count))
	sem := make(chan struct{}, concurrencyLimit)
	for _, url := range urls[:count] {
		wg.Add(1)
		url := url
		go xuiWorkflow(&wg, url.Data, ch, sem)
	}
	wg.Wait()
	close(ch)
	for proxy := range ch {
		proxies = append(proxies, proxy)
	}
	return proxies
}
func xuiWorkflow(wg *sync.WaitGroup, url string, ch chan model.Proxy, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{} // 获取一个信号量
	defer func() {
		<-sem
	}() // 释放信号量
	data := model.AddedData{
		CType: "xui",
		Data:  url,
	}
	defer database.UpdateAddedDataValidated(data)
	// 登录 - 获取cookie
	cookie := GetCookie(url)
	if cookie == "" {
		return
	}
	// 获取数据
	list := GetBoundList(url, cookie)
	// 获取inboundList
	for _, inbound := range list {

		// 返回结果
		proxy := GenDBProxy(url, inbound)
		ch <- proxy
	}
}
func GetCookie(url string) string {
	params := []byte(`username=admin&password=admin`)
	reader := bytes.NewReader(params)
	url = url + "/login"
	url = strings.ReplaceAll(url, "//", "/")
	url = strings.ReplaceAll(url, ":/", "://")
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return ""
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.76")
	//client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	client, _ := utils.GetClient(true)
	response, err := client.Do(req)
	if err != nil {
		return ""
	}
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	if string(result) == `{"success":true,"msg":"登录成功","obj":null}` {
		return response.Header.Get("Set-Cookie")
	} else {
		return ""
	}
}
func GetBoundList(url, cookie string) []model.Inbound {
	url = url + "/xui/inbound/list"
	url = strings.ReplaceAll(url, "//", "/")
	url = strings.ReplaceAll(url, ":/", "://")
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.76")
	req.Header.Add("Cookie", cookie)
	if err != nil {
		//fmt.Print(".")
		return nil
	}
	client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	msg := model.InboundMsg{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return nil
	}
	return msg.Obj
}

// GenLink 生成数据
func GenLink(url string, inbound model.Inbound) string {
	addr := SimpleRegex(url, "http(s)?://(.*?):")
	switch inbound.Protocol {
	case model.VMESS:
		{
			host := SimpleRegex(inbound.StreamSettings, `"security":"(.*?)",`)
			if host == "none" {
				host = ""
			}
			config := model.VmessConfig{
				V:    "2",
				Ps:   inbound.Remark,
				Add:  addr,
				Port: strconv.Itoa(inbound.Port),
				Id:   SimpleRegex(inbound.Settings, "\"id\":( )?\"(.*?)\","), // 请替换为你的UUID
				Aid:  SimpleRegex(inbound.Settings, "\"alterId\":( )?(.*?)\\n"),
				Net:  SimpleRegex(inbound.StreamSettings, "\"network\":( )?\"(.*?)\","),
				Type: SimpleRegex(inbound.StreamSettings, "\"type\":( )?(.*?)\\n"),
				Host: host,
				Path: SimpleRegex(inbound.StreamSettings, "\"path\":( )?\"(.*?)\","),
				Tls:  "none",
			}
			configJson, err := json.Marshal(config)
			if err != nil {
				return ""
			}
			result := "vmess://" + base64URLEncode(configJson)
			return result
		}
	case model.SOCKS:
		{
			result :=
				SimpleRegex(inbound.Settings, "\"user\":( )?\"(.*?)\"") + " " +
					SimpleRegex(inbound.Settings, "\"pass\":( )?\"(.*?)\"")
			return result
		}
	case model.HTTP:
		{
			result :=
				SimpleRegex(inbound.Settings, "\"user\":( )?\"(.*?)\"") + " " +
					SimpleRegex(inbound.Settings, "\"pass\":( )?\"(.*?)\"")
			return result
		}
	}
	return ""
}
func GenDBProxy(url string, inbound model.Inbound) model.Proxy {
	proxy := model.Proxy{
		CType:             "xui",
		Protocol:          string(inbound.Protocol),
		Host:              SimpleRegex(url, "http(s)?://(.*?):"),
		Port:              inbound.Port,
		Validated:         false,
		Latency:           99999,
		ValidateDate:      time.Now(),
		ToValidateDate:    time.Now().Add(24 * time.Hour),
		ValidateFailedCnt: 0,
		User:              "",
		Pass:              "",
		Link:              "",
		Country:           GetIPCountry(url),
	}
	switch inbound.Protocol {
	case model.VMESS:
		{
			proxy.Link = GenLink(url, inbound)
			break
		}
	case model.HTTP:
		{
			proxy.User = SimpleRegex(inbound.Settings, "\"user\":( )?\"(.*?)\"")
			proxy.Pass = SimpleRegex(inbound.Settings, "\"pass\":( )?\"(.*?)\"")
			proxy.Link = fmt.Sprintf("%s://%s:%s@%s:%d", proxy.Protocol, proxy.User, proxy.Pass, proxy.Host, proxy.Port)
			break
		}
	case model.SOCKS:
		{
			proxy.User = SimpleRegex(inbound.Settings, "\"user\":( )?\"(.*?)\"")
			proxy.Pass = SimpleRegex(inbound.Settings, "\"pass\":( )?\"(.*?)\"")
			proxy.Link = fmt.Sprintf("%s://%s:%s@%s:%d", proxy.Protocol, proxy.User, proxy.Pass, proxy.Host, proxy.Port)
			break
		}
	default:
		proxy.User = SimpleRegex(inbound.Settings, "\"user\":( )?\"(.*?)\"")
		proxy.Pass = SimpleRegex(inbound.Settings, "\"pass\":( )?\"(.*?)\"")
		proxy.Link = fmt.Sprintf("%s://%s:%s@%s:%d", proxy.Protocol, proxy.User, proxy.Pass, proxy.Host, proxy.Port)
		break
	}
	proxy.Link = strings.ReplaceAll(proxy.Link, "://:@", "://")
	return proxy
}
func SimpleRegex(text, regex string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(text)
	if len(match) > 0 {
		return match[len(match)-1]
	}
	return ""
}
func base64URLEncode(data []byte) string {
	encoding := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	var result string
	var val int
	var valb int
	for i := 0; i < len(data); i++ {
		val = (val << 8) | int(data[i])
		valb += 8
		for valb >= 6 {
			result += string(encoding[(val>>uint(valb-6))&0x3F])
			valb -= 6
		}
	}
	if valb > 0 {
		result += string(encoding[(val<<uint(6-valb))&0x3F])
	}
	return result
}

// GetIPCountry 获取国家
func GetIPCountry(url string) string {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Get("https://get.geojs.io/v1/ip/country/" + SimpleRegex(url, "http(s)?://(.*?):"))
	if err != nil {
		fmt.Println("GetCountryError " + err.Error())
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("GetCountryError " + err.Error())
		return ""
	}
	country := string(body)
	if len(country) > 5 {
		return ""
	}
	return strings.ReplaceAll(country, "\n", "")
}

package validator

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"proxyCrawler/model"
	"strconv"
	"time"
)

func checkHTTPProxy(proxyURL, username, password string) int64 {
	prx, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Println("Invalid prx URL:", err)
		return 0
	}

	prx.User = url.UserPassword(username, password)

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(prx),
		},
		Timeout: 5 * time.Second,
	}
	start := time.Now()
	resp, err := client.Get("http://www.baidu.com")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	// 验证响应头
	elapsed := time.Since(start).Milliseconds()
	return elapsed
}
func checkSocksProxy(proxyAddr, username, password string) int64 {
	auth := &proxy.Auth{
		User:     username,
		Password: password,
	}

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, auth, proxy.Direct)
	if err != nil {
		fmt.Println("Error creating SOCKS5 dialer:", err)
		return 0
	}

	transport := &http.Transport{}
	transport.Dial = dialer.Dial

	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
	start := time.Now()
	resp, err := client.Get("http://www.baidu.com")
	// 验证响应头
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	elapsed := time.Since(start).Milliseconds()
	return elapsed
}
func checkValid(proxy model.Proxy) int64 {
	switch proxy.Protocol {
	case "http":
		{
			return checkHTTPProxy("http://"+proxy.Host+":"+strconv.Itoa(proxy.Port), proxy.User, proxy.Pass)
		}
	case "https":
		{
			return checkHTTPProxy("https://"+proxy.Host+":"+strconv.Itoa(proxy.Port), proxy.User, proxy.Pass)
		}
	case "socks":
		{
			return checkSocksProxy(proxy.Host+":"+strconv.Itoa(proxy.Port), proxy.User, proxy.Pass)
		}
	default:
		return 0
	}
}

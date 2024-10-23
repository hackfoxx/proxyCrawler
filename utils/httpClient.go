package utils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"proxyCrawler/database"
	"proxyCrawler/model"
	"proxyCrawler/utils/logger"
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
func GetClient(useProxy bool) (*http.Client, bool) {
	defaultClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, Timeout: 5 * time.Second}
	if !useProxy {
		return defaultClient, false
	}
	var dbProxy model.Proxy
	for i := 0; i < 3; i++ {
		tmp := database.GetRandomProxy(map[string]interface{}{"protocol": "http", "validated": true})
		if checkHTTPProxy("http://"+tmp.Host+":"+strconv.Itoa(tmp.Port), tmp.User, tmp.Pass) > int64(0) {
			dbProxy = database.GetRandomProxy(map[string]interface{}{"protocol": "http", "validated": true})
			break
		}
	}
	if dbProxy.Host == "" {
		return defaultClient, false
	}
	prx, err := url.Parse(dbProxy.Protocol + "://" + dbProxy.Host + ":" + strconv.Itoa(dbProxy.Port))
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Invalid prx URL: %s", err))
		return defaultClient, false
	}
	prx.User = url.UserPassword(dbProxy.User, dbProxy.Pass)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(prx),
		},
		Timeout: 5 * time.Second,
	}
	return client, true
}
func RawRequest(rawRequest string, useProxy bool) string {
	buf := bytes.NewBufferString(rawRequest)
	reader := bufio.NewReader(buf)
	req, err := http.ReadRequest(reader)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Error reading request: %s", err))
		return ""
	}

	// 修正 RequestURI 为空的问题
	req.RequestURI = ""

	// 构建完整的 URL
	req.URL, err = url.Parse("http://" + req.Host + req.URL.String())
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Error parsing URL: %s", err))
		return ""
	}

	// 发送请求
	client, _ := GetClient(useProxy)

	resp, err := client.Do(req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Error sending request: %s", err))
		return ""
	}
	defer resp.Body.Close()

	// 打印响应状态和头信息
	logger.GetLogger().Info(fmt.Sprintf("Response Status: %s", resp.Status))
	logger.GetLogger().Info(fmt.Sprintf("Response Headers: %s", resp.Header))
	// 读取响应主体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}
	return string(body)
}

package adder

import (
	"fmt"
	"io"
	"net/http"
	"proxyCrawler/database"
	"proxyCrawler/model"
	"proxyCrawler/module/crawler"
	"proxyCrawler/utils"
	"regexp"
)

func URLAdder(CType, url string) model.DBResult {
	result := model.DBResult{}
	errResult := model.DBResult{Sum: -1}
	validType := false
	for _, s := range crawler.GetCrawlerTypes() {
		if s == CType {
			validType = true
			break
		}
	}
	if !validType {
		// 不支持的type
		return errResult
	}
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		fmt.Println("GetUrlsError" + err.Error())
		return result
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("resp Error" + err.Error())
		return errResult
	}
	rawUrls := string(body)
	pattern := ""
	if CType == "xui" {
		pattern = `https?://.*?:54321`
	} else if CType == "local" {
		pattern = "(http|https|socks):(.*:.*@)?(.*?):\\d*|vmess://.*"
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return errResult
	}
	matches := re.FindAllString(rawUrls, -1)
	if matches == nil {
		fmt.Println("No matches found")
	}
	var addedData []model.AddedData
	for i, match := range matches {
		if i%1000 == 0 {
			result = utils.AddDBResult(result, database.SetAddedDataList(addedData, false))
			addedData = nil
			fmt.Printf("已上传: %d / %d\n", result.Sum, len(matches))
		}
		addedData = append(addedData, model.AddedData{
			CType:     CType,
			Data:      match,
			Validated: false,
		})
	}
	result = utils.AddDBResult(result, database.SetAddedDataList(addedData, false))
	fmt.Printf("上传完成: %s", utils.ReadDBResult("", result))
	return result
}

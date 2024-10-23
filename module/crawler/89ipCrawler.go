package crawler

import (
	"fmt"
	"io"
	"net/http"
	"proxyCrawler/model"
	"proxyCrawler/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type c89ipCrawler struct {
}

func (f c89ipCrawler) Run() []model.Proxy {
	var result []model.Proxy
	count := getSum()
	url := "http://api.89ip.cn/tqdl.html?api=1&num=" + count
	//client := utils.GetClient()
	client := http.DefaultClient
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return result
	}
	re := regexp.MustCompile("(\\d{1,3}\\.){1,3}\\d{1,3}:\\d{1,5}")
	matched := re.FindAllString(string(bytes), -1)
	for _, s := range matched {
		split := strings.Split(s, ":")
		port, _ := strconv.Atoi(split[1])
		result = append(result, model.Proxy{
			CType:             "89ip",
			Protocol:          "http",
			Host:              split[0],
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
		})
	}
	return result
}

func getSum() string {
	client, _ := utils.GetClient(false)
	resp1, err := client.Get("https://www.89ip.cn")
	if err != nil {
		fmt.Println(err)
		return "5000"
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		fmt.Println(err)
		return "5000"
	}
	sum := SimpleRegex(string(body1), "IP总量：(.*?)个")
	if sum == "0" {
		sum = "5000"
	}
	return sum
}

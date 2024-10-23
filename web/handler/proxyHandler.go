package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"proxyCrawler/database"
	"proxyCrawler/jobs"
	"proxyCrawler/model"
	"proxyCrawler/module/adder"
	"proxyCrawler/utils"
	"strconv"
)

func SetHandler(r *gin.Engine) *gin.Engine {
	v1 := r.Group("/v1")
	v1.GET("/proxies", proxies)
	v1.GET("/proxy", random)
	v1.GET("/add", addUrlProxies)
	v1.POST("/add", addRawProxies)
	v1.GET("/crawler", getCrawler)
	v1.GET("/start", startCrawler)
	v1.GET("/tp", getCount)
	return r
}

func random(c *gin.Context) {
	conditions := genConditions(c)
	if c.Query("validated") == "1" || c.Query("validated") == "true" {
		conditions["validated"] = true
	}
	if c.Query("r") == "link" {
		c.JSON(http.StatusOK, database.GetRandomProxy(conditions).Link)
	} else {
		c.JSON(http.StatusOK, database.GetRandomProxy(conditions))
	}
}

func addUrlProxies(c *gin.Context) {
	fetcherType := c.Query("c_type")
	link := c.Query("link")
	count := adder.URLAdder(fetcherType, link).Added
	c.JSON(http.StatusOK, "成功添加"+strconv.FormatInt(count, 10)+"个代理")
	jobs.CrawlerJob(fetcherType)
}

// proxies /
func proxies(c *gin.Context) {
	var result model.HttpResult
	conditions := genConditions(c)
	proxies := database.GetProxies(conditions)
	result.StatusCode = http.StatusOK
	result.Msg = "本次共获取到" + strconv.Itoa(len(proxies)) + "个代理"
	if c.Query("t") == "l" {
		var links []string
		for _, proxy := range proxies {
			links = append(links, proxy.Link)
		}
		result.Object = links
	} else {
		result.Object = proxies
	}
	c.JSON(http.StatusOK, result)
}

/*
添加代理
@params

	link="http://example.com/urls.txt" 需要添加的文件 #需要强鉴权
	ftype=(local|xui|ip3366) fetcherType 扫描器类型
*/
func addRawProxies(c *gin.Context) {
	var tp struct {
		CType   string   `json:"type"`
		Proxies []string `json:"proxies"`
	}
	// 如果 Content-Type 是 application/json，则 BindJSON 会被调用
	if c.ContentType() == "application/json" {
		if err := c.BindJSON(&tp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 否则，尝试从表单数据中获取
		if err := c.Bind(&tp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	fmt.Println(tp.CType)
	fmt.Println(tp.Proxies)
	result := adder.RawAdder(tp.CType, tp.Proxies)
	c.String(http.StatusOK, "%s", utils.ReadDBResult("添加结果\n", result))
}

// 统计代理数量
func getCount(c *gin.Context) {
	count := database.GetProxyCount()
	c.JSON(http.StatusOK, count)
}

func genConditions(c *gin.Context) map[string]interface{} {
	conditions := make(map[string]interface{})
	// web参数对应数据库字段名
	paramMapping := map[string]string{
		"c": "country",
		"p": "protocol",
		"s": "c_type",
		//"v":"validated",
		//"t":link
	}
	for k, m := range paramMapping {
		v := c.Query(k)
		if v != "" {
			conditions[m] = v
		}
	}
	if c.Query("v") == "1" || c.Query("v") == "true" {
		conditions["validated"] = true
	}
	return conditions
}

/*
返回clash配置文件
*/
func clash(c *gin.Context) {
	c.JSON(http.StatusOK, "正在开发中...")
}
func getCrawler(c *gin.Context) {
	crawlers := database.GetCrawlers(nil)
	c.JSON(http.StatusOK, crawlers)
}
func startCrawler(c *gin.Context) {
	cType := c.Param("type")
	job := jobs.CrawlerJob(cType)
	c.JSON(http.StatusOK, job)
}

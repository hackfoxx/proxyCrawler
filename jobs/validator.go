package jobs

import (
	"fmt"
	"proxyCrawler/database"
	"proxyCrawler/module/validator"
	"proxyCrawler/utils"
	"proxyCrawler/utils/logger"
)

var isRunning bool

func ValidatorJob() {
	if isRunning {
		logger.GetLogger().Info("验证器正在运行...")
		return
	}
	if database.GetCrawlers(map[string]interface{}{"is_running": true}) == nil {
		logger.GetLogger().Info("有扫描器正在运行，本次跳过验证！")
		return
	}
	isRunning = true
	defer func() { isRunning = false }()
	logger.GetLogger().Info("开始验证...")
	proxies := database.GetToValidProxies()
	logger.GetLogger().Info(fmt.Sprintf("本次需要验证 %d 个代理\n", len(proxies)))
	dbResult := validator.ValidProxies(proxies)
	logger.GetLogger().Info(utils.ReadDBResult("验证结果", dbResult))
}

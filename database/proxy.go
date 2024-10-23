package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"proxyCrawler/model"
	"proxyCrawler/utils/logger"
	"time"
)

func SetProxy(proxy model.Proxy, doUpdate bool) model.DBResult {
	prx := proxy
	var result = model.DBResult{
		Added:   0,
		Updated: 0,
		Error:   0,
		Deleted: 0,
		Sum:     1,
	}
	tx := db.Where("host = ? and port=?", proxy.Host, proxy.Port).First(&prx)
	if tx.Error != nil { //如果报错
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) { // 如果未找到代理
			// 尝试插入
			tx = db.Create(&proxy)
			if tx.Error != nil {
				fmt.Println(tx.Error)
				result.Error++
			}
			result.Added++
		} else { //其它报错
			fmt.Println(tx.Error)
			result.Error++
		}
	} else { // 如果没报错
		if doUpdate { // 执行更新
			// 尝试更新
			tx = db.Model(&proxy).Select("*").Updates(proxy)
			if tx.Error != nil {
				fmt.Println(tx.Error)
				result.Error++
			} else {
				result.Updated++
			}
		} else { // 执行跳过
			result.Continue++
		}
	}
	return result
}

// SetProxies 如果不存在则插入，如果存在则 doUpdate 如存在是否更新
func SetProxies(proxies []model.Proxy, doUpdate bool) model.DBResult {
	var result = model.DBResult{
		Added:   0,
		Updated: 0,
		Error:   0,
		Deleted: 0,
		Sum:     int64(len(proxies)),
	}
	for _, proxy := range proxies {
		prx := proxy
		// 尝试按host和port查找代理
		tx := db.Where("host = ? and port=?", proxy.Host, proxy.Port).First(&prx)
		if tx.Error != nil { //如果报错
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) { // 如果未找到代理
				// 尝试插入
				tx = db.Create(&proxy)
				if tx.Error != nil {
					logger.GetLogger().Error(tx.Error.Error())
					result.Error++
				}
				result.Added++
			} else { //其它报错
				result.Error++
				logger.GetLogger().Error(tx.Error.Error())
			}
		} else { // 如果没报错
			if doUpdate { // 执行更新
				// 尝试更新
				tx = db.Model(&proxy).Select("*").Updates(proxy)
				if tx.Error != nil {
					logger.GetLogger().Error(tx.Error.Error())
					result.Error++
				} else {
					result.Updated++
				}
			} else { // 执行跳过
				result.Continue++
			}
		}
	}
	return result
}

func GetProxies(conditions map[string]interface{}) []model.Proxy {
	var proxies []model.Proxy
	tx := db.Where(conditions).Find(&proxies)
	if tx.Error != nil {
		logger.GetLogger().Error(tx.Error.Error())
		return nil
	}
	return proxies
}
func GetToValidProxies() []model.Proxy {
	var proxies []model.Proxy
	tx := db.Where("latency=? or to_validate_date < ?", 0, time.Now()).Find(&proxies)
	if tx.Error != nil {
		logger.GetLogger().Error(tx.Error.Error())
		return nil
	}
	return proxies
}
func DeleteProxies(proxies []model.Proxy) model.DBResult {
	tx := db.Delete(proxies)
	result := model.DBResult{
		Added:   0,
		Updated: 0,
		Deleted: tx.RowsAffected,
		Error:   0,
		Sum:     int64(len(proxies)),
	}
	result.Error = result.Sum - result.Deleted
	return result
}
func GetRandomProxy(conditions map[string]interface{}) model.Proxy {
	proxy := model.Proxy{}
	db.Order("RAND()").First(&proxy, conditions)
	return proxy
}
func SaveProxy(proxy model.Proxy) bool {
	tx := db.Model(&proxy).Select("*").Updates(proxy)
	return tx.Error == nil
}
func GetProxyCount() model.CountResponse {
	var count = model.CountResponse{}
	db.Model(&model.Proxy{}).Where("validated=1").Count(&count.Validated)
	db.Model(&model.Proxy{}).Count(&count.Total)
	db.Model(&model.Proxy{}).Where("validated=1 and protocol like 'http%'").Count(&count.ValidHttp)
	db.Model(&model.Proxy{}).Where("protocol like 'http%'").Count(&count.TotalHttp)

	db.Model(&model.Proxy{}).Where("validated=1 and protocol like 'socks%'").Count(&count.ValidSocks)
	db.Model(&model.Proxy{}).Where("protocol like 'socks%'").Count(&count.TotalSocks)

	db.Model(&model.Proxy{}).Where("validated=1 and protocol = 'vmess'").Count(&count.ValidVmess)
	db.Model(&model.Proxy{}).Where("protocol = 'vmess'").Count(&count.TotalVmess)
	return count
}

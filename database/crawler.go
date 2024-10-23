package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"proxyCrawler/model"
)

// SetCrawlers 如果不存在则插入，如果存在则更新
func SetCrawlers(crawlers []model.Crawler, doUpdate bool) model.DBResult {
	crs := crawlers
	var result = model.DBResult{
		Added:   0,
		Updated: 0,
		Error:   0,
		Deleted: 0,
		Sum:     int64(len(crawlers)),
	}
	for _, crawler := range crawlers {
		// 尝试按host和port查找代理
		tx := db.Where("type=?", crawler.Type).First(&crs)
		if tx.Error != nil {
			// 如果未找到代理
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				// 尝试插入
				tx = db.Create(&crawler)
				if tx.Error != nil {
					result.Error++
				}
				result.Added++
			} else { //如果代理已存在
				result.Error++
			}
		}
		if doUpdate {
			// 尝试更新
			tx = db.Model(&crawler).Select("*").Updates(crawler)
			//tx = db.Updates(&crawler)
			if tx.Error != nil {
				result.Error++
			} else {
				result.Updated++
			}
		} else {
			result.Continue++
		}
	}
	return result
}

func GetCrawlers(conditions map[string]interface{}) []model.Crawler {
	var crawlers []model.Crawler
	tx := db.Where(conditions).Find(&crawlers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return nil
	}
	return crawlers
}
func DeleteCrawlers(crawlers []model.Crawler) model.DBResult {
	tx := db.Delete(crawlers)
	result := model.DBResult{
		Added:   0,
		Updated: 0,
		Deleted: tx.RowsAffected,
		Error:   0,
		Sum:     int64(len(crawlers)),
	}
	result.Error = result.Sum - result.Deleted
	return result
}
func GetCrawler(cType string) model.Crawler {
	var crawler model.Crawler
	db.Where("type=?", cType).First(&crawler)
	return crawler
}

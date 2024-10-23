package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"proxyCrawler/model"
)

// SetAddedData 返回值{-1：报错，1：添加，2：更新，3：跳过}
func SetAddedData(data model.AddedData, doUpdate bool) model.DBResult {
	dta := data
	var result = model.DBResult{
		Added:   0,
		Updated: 0,
		Deleted: 0,
		Error:   0,
		Sum:     1,
	}
	// 尝试按host和port查找代理
	tx := db.Where("c_type=? and data=?", data.CType, data.Data).First(&dta)
	if tx.Error != nil {
		// 如果未找到代理
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			// 尝试插入
			tx = db.Create(&data)
			if tx.Error != nil {
				result.Error++
			}
			result.Added++
		} else { // 其它报错
			result.Error++
		}
	} else if doUpdate {
		//如果代理已存在
		// 尝试更新
		tx = db.Model(&data).Select("*").Updates(data)
		if tx.Error != nil {
			fmt.Println(tx.Error)
			result.Error++
		} else {
			result.Updated++
		}
	} else {
		result.Continue++
	}
	return result
}

// SetAddedDataList 如果不存在则插入，如果存在则更新 continueUpdate 如果为false则强制更新，如为true则不更新
func SetAddedDataList(data []model.AddedData, doUpdate bool) model.DBResult {
	var result = model.DBResult{
		Added:   0,
		Updated: 0,
		Deleted: 0,
		Error:   0,
		Sum:     int64(len(data)),
	}
	for _, line := range data {
		// 尝试按host和port查找代理
		tx := db.Where("c_type=? and data=?", line.CType, line.Data).First(&data)
		if tx.Error != nil {
			// 如果未找到代理
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				// 尝试插入
				tx = db.Create(&line)
				if tx.Error != nil {
					result.Error++
					continue
				}
				result.Added++
				continue
			} else { // 其它报错
				result.Error++
				continue
			}
		} else if doUpdate {
			//如果代理已存在
			// 尝试更新
			tx = db.Model(&line).Select("*").Updates(line)
			if tx.Error != nil {
				fmt.Println(tx.Error)
				result.Error++
				continue
			} else {
				result.Updated++
				continue
			}
		} else {
			result.Continue++
		}
	}
	return result
}
func GetAddedData(conditions map[string]interface{}) []model.AddedData {
	var data []model.AddedData
	tx := db.Where(conditions).Find(&data)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return nil
	}
	return data
}
func DeleteAddedData(data []model.AddedData) model.DBResult {
	tx := db.Delete(data)
	result := model.DBResult{
		Added:   0,
		Updated: 0,
		Deleted: tx.RowsAffected,
		Error:   0,
		Sum:     int64(len(data)),
	}
	result.Error = result.Sum - result.Deleted
	return result
}
func UpdateAddedDataValidated(data model.AddedData) bool {
	tx := db.Where(&data).Updates(model.AddedData{Validated: true})
	return tx.Error == nil
}

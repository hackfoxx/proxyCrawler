package adder

import (
	"fmt"
	"proxyCrawler/database"
	"proxyCrawler/model"
	"proxyCrawler/utils"
)

func RawAdder(cType string, proxies []string) model.DBResult {
	result := model.DBResult{}
	var tmp []model.AddedData
	for i, proxy := range proxies {
		if i%1000 == 0 {
			result = utils.AddDBResult(result, database.SetAddedDataList(tmp, false))
			tmp = nil
			fmt.Printf("上传进度: %d / %d\n", result.Sum, len(proxies))
		}
		if checkUrl(cType, proxy) {
			tmp = append(tmp, model.AddedData{
				CType:     cType,
				Data:      proxy,
				Validated: false,
			})
		} else {
			result.Error++
		}
	}
	result = utils.AddDBResult(result, database.SetAddedDataList(tmp, false))
	fmt.Printf("上传进度: %d / %d\n", result.Sum, len(proxies))
	return result
}

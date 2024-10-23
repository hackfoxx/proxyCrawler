package utils

import (
	"proxyCrawler/model"
	"strconv"
)

func ReadDBResult(title string, result model.DBResult) string {
	msg := title
	msg += "新增: " + strconv.Itoa(int(result.Added)) + " "
	msg += "更新: " + strconv.Itoa(int(result.Updated)) + " "
	msg += "删除: " + strconv.Itoa(int(result.Deleted)) + " "
	msg += "跳过: " + strconv.Itoa(int(result.Continue)) + " "
	msg += "报错: " + strconv.Itoa(int(result.Error)) + " "
	msg += "总数: " + strconv.Itoa(int(result.Sum))
	return msg
}
func AddDBResult(r1, r2 model.DBResult) model.DBResult {
	r1.Added += r2.Added
	r1.Updated += r2.Updated
	r1.Deleted += r2.Deleted
	r1.Continue += r2.Continue
	r1.Error += r2.Error
	r1.Sum += r2.Sum
	return r1
}

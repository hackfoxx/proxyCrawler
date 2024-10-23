package adder

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"proxyCrawler/database"
	"proxyCrawler/model"
	"proxyCrawler/utils"
	"regexp"
)

func checkUrl(cType, url string) bool {
	matched, _ := regexp.MatchString(`(https?|socks)://((\w|\d|-)*?:(\w|\d-)*?@)?(\d|\w|\.|-)*?:\d{1,5}`, url)
	switch cType {
	case "xui":
		{
			matched, _ = regexp.MatchString(`https?://(.*):54321/?`, url)
		}
	case "local":
		{
			//matched, _ = regexp.MatchString(`(https?|socks)://((\w|\d|-)*?:(\w|\d-)*?@)?(\d|\w|\.|-)*?:\d{1,5}`, url)
			break
		}
	}
	return matched
}

func FileAdder(cType, fileName string) model.DBResult {
	var result model.DBResult
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		url := string(line)
		fmt.Println(url)
		if checkUrl(cType, url) {
			result = utils.AddDBResult(result, database.SetAddedData(model.AddedData{
				CType:     cType,
				Data:      url,
				Validated: false,
			}, false))
		}
	}
	return result
}

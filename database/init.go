package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"proxyCrawler/config"
	"proxyCrawler/model"
)

var (
	db          = Init()
	dbAddedData = model.AddedData{}
	dbProxy     = model.Proxy{}
	dbCrawler   = model.Crawler{}
)

func Init() *gorm.DB {
	dbConfig := config.GetConfig().Database
	if dbConfig.Host == "" {
		fmt.Println("数据库配置文件读取失败")
		os.Exit(-1)
	}
	username := dbConfig.Username
	password := dbConfig.Password
	host := dbConfig.Host
	port := dbConfig.Port
	Dbname := dbConfig.Dbname
	timeout := dbConfig.Timeout
	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		fmt.Println("数据库连接失败: " + err.Error())
		os.Exit(-1)
	}
	err = db.AutoMigrate(&dbAddedData, &dbProxy, &dbCrawler)
	if err != nil {
		fmt.Println("数据表创建失败" + err.Error())
		os.Exit(-1)
	}
	return db
}

package config

import (
	"github.com/spf13/viper"
	"log"
)

var cfg Config

type Config struct {
	Web struct {
		Addr          string
		BasePath      string
		Authorization string
	}
	Database struct {
		Username string
		Password string
		Host     string
		Port     int
		Dbname   string
		Timeout  string
	}
}

func GetConfig() Config {
	viper.SetConfigName("config") // 配置文件名称 (不带扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的扩展名不是 yaml，可以设置这个
	viper.AddConfigPath(".")      // 添加配置文件所在的路径
	// 设置默认值
	viper.SetDefault("web.addr", ":8080")
	viper.SetDefault("web.base_path", "/v1")
	viper.SetDefault("web.authorization", "a57f1abe-c6df-4a9b-82ad-a29cf1304399")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return cfg
}

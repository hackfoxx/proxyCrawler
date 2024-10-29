package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var Cfg *Config
var once sync.Once
var defaultConfig = getDefault()

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

func Init(fp string) {
	once.Do(func() {
		LoadConfig(fp)
	})
}
func LoadConfig(fp string) {
	var conf Config
	fmt.Println("正在使用配置文件[", fp, "]")
	viper.SetConfigFile(fp)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("配置文件[%s]不存在，请选择下一步操作\n1. 在[%s]生成默认配置\n2. 生成默认[./config.yml]配置\n0. 退出\n", fp, fp)
		var selection int
		fmt.Scanf("%d", &selection)
		switch selection {
		case 1:
			err := saveConfig(defaultConfig, fp)
			if err != nil {
				log.Fatalf("%s", err.Error())
			} else {
				fmt.Printf("已生成配置文件[%s]\n", fp)
			}
			os.Exit(0)
		case 2:
			err := saveConfig(defaultConfig, "./config.yml")
			if err != nil {
				log.Fatalf("%s", err.Error())
			} else {
				fmt.Printf("已生成配置文件[%s]\n", "./config.yml")
			}
			os.Exit(0)
		case 0:
			os.Exit(0)
		default:
			os.Exit(-1)
		}
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	setDefault(&conf)
	Cfg = &conf
}
func Check(fp string) bool {
	viper.SetConfigFile(fp)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("配置文件[%s]不存在: %v\n", fp, err)
		return false
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("配置文件[%s]解析失败: %v\n", fp, err)
		return false
	}
	return true
}

// saveConfig 生成默认配置文件
func saveConfig(cfg *Config, filename string) error {
	data, err := yaml.Marshal(cfg) // 将结构体序列化为 YAML 格式
	if err != nil {
		return fmt.Errorf("unable to marshal config to YAML: %w", err)
	}

	// 将 YAML 数据写入文件
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("unable to write config to file: %w", err)
	}
	return nil
}
func SaveDefault(fp string) error {
	return saveConfig(defaultConfig, fp)
}

// setDefault 装配默认值
func setDefault(cfg *Config) {
	if cfg.Web.Addr == "" {
		cfg.Web.Addr = defaultConfig.Web.Addr
	}
	if cfg.Web.BasePath == "" {
		cfg.Web.BasePath = defaultConfig.Web.BasePath
	}
	if cfg.Web.Authorization == "" {
		cfg.Web.Authorization = defaultConfig.Web.Authorization
	}
	if cfg.Database.Host == "" {
		cfg.Database.Host = defaultConfig.Database.Host
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = defaultConfig.Database.Port
	}
	if cfg.Database.Username == "" {
		cfg.Database.Username = defaultConfig.Database.Username
	}
	if cfg.Database.Password == "" {
		cfg.Database.Password = defaultConfig.Database.Password
	}
	if cfg.Database.Dbname == "" {
		cfg.Database.Dbname = defaultConfig.Database.Dbname
	}
	if cfg.Database.Timeout == "" {
		cfg.Database.Timeout = defaultConfig.Database.Timeout
	}
}

// getDefault 初始化 defaultConfig变量
func getDefault() *Config {
	var cfg Config
	cfg.Web.Addr = ":8080"
	cfg.Web.BasePath = "/v1"
	cfg.Web.Authorization = "a57f1abe-c6df-4a9b-82ad-a29cf1304399"
	cfg.Database.Host = "127.0.0.1"
	cfg.Database.Port = 3306
	cfg.Database.Username = "root"
	cfg.Database.Password = "root"
	cfg.Database.Dbname = "crawler"
	cfg.Database.Timeout = "10s"
	return &cfg
}

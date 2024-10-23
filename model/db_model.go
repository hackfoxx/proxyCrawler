package model

import "time"

// AddedData 被添加的数据
type AddedData struct {
	CType     string `gorm:"primaryKey"`
	Data      string `gorm:"primaryKey"`
	Validated bool
}

// Proxy -------存入数据库的代理格式----------
type Proxy struct {
	CType             string    // 爬取器名称
	Protocol          string    // 协议
	Host              string    `gorm:"primaryKey"` // 主机名(IP地址或域名)
	Port              int       `gorm:"primaryKey"` // 端口
	Validated         bool      // 是否有效
	Latency           int       // 延迟
	ValidateDate      time.Time // 上次验证的时间
	ToValidateDate    time.Time // 计划下次验证的时间
	ValidateFailedCnt int       // 验证失败次数
	User              string    // 用户名
	Pass              string    // 密码
	Link              string    // Vmess链接
	Country           string    // 国家
}

// Crawler 扫描器
type Crawler struct {
	Type           string    `gorm:"primaryKey"` //扫描器类型 - 需要在fetcher/Fetcher.go#fetcherMap里面注册
	Enable         bool      // 是否启用
	SumProxiesCnt  int       // 总共获取的代理数
	LastProxiesCnt int       // 最后一次获取的代理数量
	LastFetchDate  time.Time // 最后一次执行的时间
	ToFetchDate    time.Time // 计划下次执行时间
	// todo 使用数据库来实现单例运行的效果
	IsRunning bool
}

type DBResult struct {
	Msg      string
	Added    int64 //本次新增
	Updated  int64 //本次更新
	Deleted  int64 //本次删除
	Continue int64 //已存在未更新
	Error    int64
	Sum      int64 //总数
}

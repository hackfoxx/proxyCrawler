package model

// Msg -----XUI常规响应格式--------
type Msg struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Obj     interface{} `json:"obj"`
}

// InboundMsg -------XUI 解析 Inbound数组------------
type InboundMsg struct {
	Success bool      `json:"success"`
	Msg     string    `json:"msg"`
	Obj     []Inbound `json:"obj"`
}

// Inbound ------XUI Inbound-----------
type Inbound struct {
	Id         int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	UserId     int    `json:"-"`
	Up         int64  `json:"up" form:"up"`
	Down       int64  `json:"down" form:"down"`
	Total      int64  `json:"total" form:"total"`
	Remark     string `json:"remark" form:"remark"`
	Enable     bool   `json:"enable" form:"enable"`
	ExpiryTime int64  `json:"expiryTime" form:"expiryTime"`
	// config part
	Listen         string   `json:"listen" form:"listen"`
	Port           int      `json:"port" form:"port" gorm:"unique"`
	Protocol       Protocol `json:"protocol" form:"protocol"`
	Settings       string   `json:"settings" form:"settings"`
	StreamSettings string   `json:"streamSettings" form:"streamSettings"`
	Tag            string   `json:"tag" form:"tag" gorm:"unique"`
	Sniffing       string   `json:"sniffing" form:"sniffing"`
}

// VmessConfig --------用于生成Vmess链接---------
type VmessConfig struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	Id   string `json:"id"`
	Aid  string `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	Tls  string `json:"tls"`
}

// Protocol -------XUI支持的协议类型-----------
type Protocol string

const (
	VMESS       Protocol = "vmess"
	VLESS       Protocol = "vless"
	TROJAN      Protocol = "trojan"
	SHADOWSOCKS Protocol = "shadowsocks"
	DOKODEMO    Protocol = "dokodemo-door"
	MTPROTO     Protocol = "mtproto"
	SOCKS       Protocol = "socks"
	HTTP        Protocol = "http"
)

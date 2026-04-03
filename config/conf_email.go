package config

type Email struct {
	Host     string `json:"host" yaml:"host"`         // 邮件服务器地址
	Port     int    `json:"port" yaml:"port"`         // 邮件服务器端口
	From     string `json:"from" yaml:"from"`         // 发件人邮箱地址
	Nickname string `json:"nickname" yaml:"nickname"` // 发件人昵称
	Secret   string `json:"secret" yaml:"secret"`     // 发件人密码，用于验证
	IsSSL    bool   `json:"is_ssl" yaml:"is_ssl"`     // 是否ssl加密
}

package config

type ES struct {
	URL            string `json:"url" yaml:"url"`                           // es服务url
	Username       string `json:"username" yaml:"username"`                 // 用户名
	Password       string `json:"password" yaml:"password"`                 // 密码
	IsConsolePrint bool   `json:"is_console_print" yaml:"is_console_print"` // 是否打印
}

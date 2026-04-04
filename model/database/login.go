package database

import "server/global"

// 登录信息表
type Login struct {
	global.MODEL        // 附加通用表结构
	UserID       uint   `json:"user_id"` // 用户id
	User         User   `json:"user" gorm:"foreignKey:UserID"`
	LoginMethod  string `json:"login_method"` // 登录方式
	IP           string `json:"ip"`           // IP地址
	Address      string `json:"address"`      // 登录地址
	OS           string `json:"os"`           // 操作系统
	DeviceInfo   string `json:"device_info"`  // 设备信息
	BrowserInfo  string `json:"browser_indo"` // 浏览器信息
	Status       int    `json:"status"`       // 登录状态
}

package database

import "server/global"

// 友情链接
type FriendLink struct {
	global.MODEL        // 附加通用的数据库结构
	Logo         string `json:"logo" gorm:"size:255"` // logo图标
	Image        Image  `json:"-" gorm:"foreignKey:Logo;references:URL"`
	Link         string `json:"link"`        // 链接
	Name         string `json:"name"`        // 用户名
	Description  string `json:"description"` // 描述
}

package database

import "server/global"

// 打小广告
type Advertisement struct {
	global.MODEL        // 附加通用的数据库结构
	AdImage      string `json:"ad_image" gorm:"size:255"` // 图片
	Image        Image  `json:"-" gorm:"foreignKey:AdImage;references:URL"`
	Link         string `json:"link"`    // 链接
	Title        string `json:"title"`   // 标题
	Content      string `json:"content"` // 内容
}

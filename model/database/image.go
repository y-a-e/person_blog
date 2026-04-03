package database

import (
	"server/global"
	"server/model/appTypes"
)

type Image struct {
	global.MODEL                   // 附加通用的数据库结构
	Name         string            `json:"name"`                       // 名称
	URL          string            `json:"url" gorm:"size:255;unique"` // 路径
	Category     appTypes.Category `json:"category"`                   // 类别
	Storage      appTypes.Storage  `json:"storage"`                    // 存储类型
}

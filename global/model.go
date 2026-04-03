package global

import (
	"time"

	"gorm.io/gorm"
)

// 通用表附加结构
type MODEL struct {
	ID        uint           `json:"id" gorm:"primaryKey"` // 主键
	CreatedAt time.Time      `json:"created_at"`           // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`           // 更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`       // 软删除时间
}

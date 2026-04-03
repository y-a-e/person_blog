package database

import (
	"server/global"

	"github.com/gofrs/uuid"
)

// 反馈/回复
type Feedback struct {
	global.MODEL           // 附加通用的数据库结构
	UserUUID     uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`               // 用户id
	User         User      `json:"-" gorm:"foreignKey:UserUUID;references:UUID"` // 关联用户，注:references:UUID是告知关联到是User.UUID字段，而不是主键User.ID
	Content      string    `json:"centent"`                                      // 内容
	Reply        string    `json:"reply"`                                        // 回复
}

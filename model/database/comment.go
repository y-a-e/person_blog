package database

import (
	"server/global"

	"github.com/gofrs/uuid"
)

// 评论内容
type Comment struct {
	global.MODEL           // 附加通用的数据库结构
	ArticleID    string    `json:"article_id"` // 文章ID
	PID          *uint     `json:"p_id"`       // 主评价ID
	PComment     *Comment  `json:"-" gorm:"foreignKey:PID"`
	Children     []Comment `json:"children" gorm:"foreignKey:PID"`  // 子评论
	UserUUID     uuid.UUID `sjon:"user_uuid" gorm:"type:char(36)"`  // 用户UUID
	User         User      `json:"user" gorm:"foreignKey:UserUUID;references:UUID"` // 关联的用户
	Content      string    `json:"content"`                         // 内容
}

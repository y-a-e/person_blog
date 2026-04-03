package database

import "server/global"

type ArticleLike struct {
	global.MODEL        // 附加通用的数据库结构
	ArticleID    string `json:"article_id"`                 // 文章id
	UserID       uint   `json:"user_id"`                    // 用户id
	User         User   `json:"-" gorm:"foreignKey:UserID”` // 外键
}

package database

// 文章类别
type ArticleCategory struct {
	Category string `json:"category" gorm:"primaryKey"` // 类型
	Number   int    `json:"number"`                     // 统计数量
}

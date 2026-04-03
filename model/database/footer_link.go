package database

//页脚友情链接
type FooterLink struct {
	Title string `json:"title" gorm:"primaryKey"` // 标题
	Link  string `json:"link"`                    // 链接
}

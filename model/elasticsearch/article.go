package elasticsearch

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

// 文章表
type Article struct {
	CreatedAt string `json:"created_at"` // 创建时间
	UpdateAt  string `json:"updatea_at"` // 更新时间

	Cover    string   `json:"cover"`    // 封面
	Title    string   `json:"title"`    // 标题
	Keyword  string   `json:"keyword"`  // 关键字
	Category string   `json:"category"` // 类别
	Tag      []string `json:"tags"`     // 标签
	Abstract string   `json:"abstract"` // 简介
	Content  string   `json:"content"`  // 内容

	Views    int `json:"views"`    // 浏览量
	Comments int `json:"comments"` // 评论量
	Likes    int `json:"likes"`    // 收藏量
}

// ArticleIndex 文章 ES 索引
func ArticleIndex() string {
	return "article_index"
}

// DateProperty：时间；TextProperty：文本，支持模糊查询；KeywordProperty：关键字精确匹配；IntegerNumberProperty：整数类型，用于数值计算和排序
// ArticleMapping 文章 Mapping 映射
func ArticleMapping() *types.TypeMapping {
	return &types.TypeMapping{
		Properties: map[string]types.Property{
			"created_at": types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")},
			"updated_at": types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")},
			"cover":      types.TextProperty{},
			"title":      types.TextProperty{},
			"keyword":    types.KeywordProperty{},
			"category":   types.KeywordProperty{},
			"tags":       []types.KeywordProperty{},
			"abstract":   types.TextProperty{},
			"content":    types.TextProperty{},
			"views":      types.IntegerNumberProperty{},
			"comments":   types.IntegerNumberProperty{},
			"likes":      types.IntegerNumberProperty{},
		},
	}
}

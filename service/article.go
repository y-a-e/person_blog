package service

import (
	"context"
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/elasticsearch"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"gorm.io/gorm"
)

type ArticleService struct {
}

func (articleService *ArticleService) ArticleLike(req request.ArticleLike) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var al database.ArticleLike
		var num int

		// 如果文章没收藏过就收藏，收藏过就取消收藏
		if errors.Is(tx.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&al).Error, gorm.ErrRecordNotFound) {
			if err := global.DB.Create(&database.ArticleLike{
				UserID:    req.UserID,
				ArticleID: req.ArticleID,
			}).Error; err != nil {
				return err
			}
			num = 1
		} else {
			if err := global.DB.Delete(&al).Error; err != nil {
				return err
			}
			num = -1
		}

		source := "ctx._source.likes += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).Script(&script).Do(context.TODO())
		return err
	})
}

func (articleService *ArticleService) ArticleIsLike(req request.ArticleLike) (bool, error) {
	return !errors.Is(global.DB.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&database.ArticleLike{}).Error, gorm.ErrRecordNotFound), nil
}

func (articleService *ArticleService) ArticleLikesList(info request.ArticleLikesList) (interface{}, int64, error) {
	db := global.DB.Where("user_id = ?", info.UserID)
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	// 查询ArticleLike结构数据，然后再封装成ArticleLikesList结构
	l, total, err := utils.MySQLPagination(&database.ArticleLike{}, option)
	if err != nil {
		return nil, 0, err
	}

	// 定义匿名结构体，用于返回
	var list []struct {
		Id_     string                `json:"_id"`
		Source_ elasticsearch.Article `json:"_source"`
	}

	for _, articleLike := range l {
		article, err := articleService.Get(articleLike.ArticleID)
		if err != nil {
			return nil, 0, err
		}

		article.Keyword = ""
		article.Content = ""
		article.UpdateAt = ""
		list = append(list, struct {
			Id_     string                "json:\"_id\""
			Source_ elasticsearch.Article "json:\"_source\""
		}{
			Id_:     articleLike.ArticleID,
			Source_: article,
		})
	}

	return list, total, nil
}

func (articleService *ArticleService) ArticleInfoByID(id string) (elasticsearch.Article, error) {
	go func() {
		articleView := articleService.NewArticleView()
		_ = articleView.Set(id)
	}()
	return articleService.Get(id)
}

func (articleService *ArticleService) ArticleSearch(info request.ArticleSearch) (interface{}, int64, error) {
	// 构建查询请求
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}

	// 根据查询字段查询
	if info.Query != "" {
		// Shoule 或， Must 且， Filter 过滤
		boolQuery.Should = []types.Query{
			// []types.Query >> Match map[string]MatchQuery
			// 多个查询条件 只要 title 、 keyword 、 abstract 、 content 任意一个字段 匹配，就符合条件
			{Match: map[string]types.MatchQuery{"title": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"keyword": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"abstract": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"content": {Query: info.Query}}},
		}
	}

	// 根据标签筛选 Must，精确匹配，必须包含标签
	if info.Tag != "" {
		boolQuery.Must = []types.Query{
			{Match: map[string]types.MatchQuery{"tags": {Query: info.Tag}}},
		}
	}

	// 根据类别筛选 Filter，过滤查询，只在类别当中
	if info.Category != "" {
		boolQuery.Filter = []types.Query{
			{Term: map[string]types.TermQuery{"category": {Value: info.Category}}},
		}
	}

	// 如果有查询条件，则使用 Bool 查询，否则使用 MatchAll 查询
	if boolQuery.Should != nil || boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
	}

	// 设置排序字段
	if info.Sort != "" {
		var sortField string
		switch info.Sort {
		case "time":
			sortField = "created_at"
		case "view":
			sortField = "views"
		case "comment":
			sortField = "comments"
		case "like":
			sortField = "likes"
		default:
			sortField = "created_at"
		}

		var order sortorder.SortOrder
		if info.Order != "asc" {
			order = sortorder.Desc
		} else {
			order = sortorder.Asc
		}

		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					sortField: {Order: &order},
				},
			},
		}
	}

	option := other.EsOption{
		PageInfo:       info.PageInfo,
		Index:          elasticsearch.ArticleIndex(),
		Request:        req,
		SourceIncludes: []string{"created_at", "cover", "title", "abstract", "category", "tags", "views", "comments", "likes"},
	}
	return utils.EsPagination(context.TODO(), option)
}

func (articleService *ArticleService) ArticleCategory() ([]database.ArticleCategory, error) {
	var category []database.ArticleCategory
	if err := global.DB.Find(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (articleService *ArticleService) ArticleTags() ([]database.ArticleTag, error) {
	var tag []database.ArticleTag
	if err := global.DB.Find(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (articleService *ArticleService) ArticleCreate(req request.ArticleCreate) error {
	// 重复标题
	b, err := articleService.Exits(req.Title)
	if err != nil {
		return err
	}
	if b {
		return errors.New("the article already exists")
	}

	// 定义匿名结构体，用于返回
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToCreate := elasticsearch.Article{
		CreatedAt: now,
		UpdateAt:  now,
		Cover:     req.Cover,
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tag:       req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 同时更新类别表
		if err := articleService.UpdateCategoryCount(tx, "", articleToCreate.Category); err != nil {
			return err
		}
		// 同时更新标签表
		if err := articleService.UpdateTagsCount(tx, []string{}, articleToCreate.Tag); err != nil {
			return err
		}
		// 同时更新图片类别表
		if err := utils.ChangeImagesCategory(tx, []string{articleToCreate.Cover}, appTypes.Cover); err != nil {
			return err
		}
		illustrations, err := utils.FindIllustrations(articleToCreate.Content)
		if err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, illustrations, appTypes.Illustration); err != nil {
			return err
		}
		return articleService.Create(&articleToCreate)
	})
}

func (articleService *ArticleService) ArticleDelete(req request.ArticleDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}

	// 删除文章前尝试先将其他关联的数据删除
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			articleToDelete, err := articleService.Get(id)
			if err != nil {
				return err
			}
			// 该类别计数减少
			if err := articleService.UpdateCategoryCount(tx, articleToDelete.Category, ""); err != nil {
				return err
			}
			// 该标签计数减少
			if err := articleService.UpdateTagsCount(tx, articleToDelete.Tag, []string{}); err != nil {
				return err
			}
			// 图片类型初始化
			if err := utils.InitImagesCategory(tx, []string{articleToDelete.Cover}); err != nil {
				return err
			}
			illustrations, err := utils.FindIllustrations(articleToDelete.Content)
			if err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, illustrations); err != nil {
				return err
			}
			// 同时删除该文章下的所有评论
			comments, err := ServiceGroupApp.CommentService.CommentInfoByArticleID(request.CommentInfoByArticleID{ArticleID: id})
			if err != nil {
				return err
			}
			for _, comment := range comments {
				if err := ServiceGroupApp.DeleteCommentAndChild(tx, comment.ID); err != nil {
					return err
				}
			}
		}
		return articleService.Delete(req.IDs)
	})
}

func (articleService *ArticleService) ArticleUpdate(req request.ArticleUpdate) error {
	// 定义匿名结构体，不使用Article字段更新是防止将其他未使用的字段自动初始化
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToUpdate := struct {
		UpdateAt string   `json:"updatea_at"`
		Cover    string   `json:"cover"`
		Title    string   `json:"title"`
		Keyword  string   `json:"keyword"`
		Category string   `json:"category"`
		Tag      []string `json:"tags"`
		Abstract string   `json:"abstract"`
		Content  string   `json:"content"`
	}{
		UpdateAt: now,
		Cover:    req.Cover,
		Title:    req.Title,
		Keyword:  req.Title,
		Category: req.Category,
		Tag:      req.Tags,
		Abstract: req.Abstract,
		Content:  req.Content,
	}

	// 更新文章前尝试先修改其他关联的数
	return global.DB.Transaction(func(tx *gorm.DB) error {
		oldArticle, err := articleService.Get(req.ID)
		if err != nil {
			return err
		}

		// 该类别计数修改
		if err := articleService.UpdateCategoryCount(tx, oldArticle.Category, articleToUpdate.Category); err != nil {
			return err
		}

		// 该标签计数修改
		if err := articleService.UpdateTagsCount(tx, oldArticle.Tag, articleToUpdate.Tag); err != nil {
			return err
		}

		// 封面图片类型修改
		if articleToUpdate.Cover != oldArticle.Cover {
			if err := utils.InitImagesCategory(tx, []string{oldArticle.Cover}); err != nil {
				return err
			}
			if err := utils.ChangeImagesCategory(tx, []string{articleToUpdate.Cover}, appTypes.Cover); err != nil {
				return err
			}
		}

		//文章内容图片类型修改
		oldIllustrations, err := utils.FindIllustrations(oldArticle.Content)
		if err != nil {
			return err
		}
		newIllustrations, err := utils.FindIllustrations(articleToUpdate.Content)
		if err != nil {
			return err
		}
		addedIllustrations, removedIllustrations := utils.DiffArrays(oldIllustrations, newIllustrations)
		if err := utils.InitImagesCategory(tx, removedIllustrations); err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, addedIllustrations, appTypes.Illustration); err != nil {
			return err
		}

		return articleService.Update(req.ID, articleToUpdate)
	})
}

func (articleService *ArticleService) ArticleList(info request.ArticleList) (list interface{}, total int64, err error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}

	// 根据标题查询 Must，精确匹配，必须包含标签
	if info.Title != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"title": {Query: *info.Title}}})
	}

	// 根据简介查询 Must，精确匹配，必须包含标签
	if info.Abstract != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"abstract": {Query: *info.Abstract}}})
	}

	// 根据类别筛选 Filter 过滤查询，只在类别当中
	if info.Category != nil {
		boolQuery.Filter = []types.Query{
			{
				Term: map[string]types.TermQuery{
					"category": {Value: info.Category},
				},
			},
		}
	}

	// 根据条件执行查询
	if boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"created_at": {Order: &sortorder.Desc},
				},
			},
		}
	}

	option := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  req,
	}
	return utils.EsPagination(context.TODO(), option)
}

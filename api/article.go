package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleApi struct {
}

// 根据id查询文章详情
func (articleApi *ArticleApi) ArticleInfoByID(c *gin.Context) {
	var req request.ArticleInfoByID
	err := c.ShouldBindUri(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	article, err := articleService.ArticleInfoByID(req.ID)
	if err != nil {
		global.Log.Error("failed to articleInfo : ", zap.Error(err))
		response.FailWithMessage("failed to articleInfo", c)
	}
	response.OkWithData(article, c)
}

// 根据查询字段查询文章
func (articleApi *ArticleApi) ArticleSearch(c *gin.Context) {
	var info request.ArticleSearch
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := articleService.ArticleSearch(info)
	if err != nil {
		global.Log.Error("failed to articleSearch : ", zap.Error(err))
		response.FailWithMessage("failed to articleSearch", c)
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// 获取所有文章类别及数量
func (articleApi *ArticleApi) ArticleCategory(c *gin.Context) {
	category, err := articleService.ArticleCategory()
	if err != nil {
		global.Log.Error("failed to articleCategory : ", zap.Error(err))
		response.FailWithMessage("failed to articleCategory", c)
	}
	response.OkWithData(category, c)
}

// 获取所有文章标签及数量
func (articleApi *ArticleApi) ArticleTags(c *gin.Context) {
	tags, err := articleService.ArticleTags()
	if err != nil {
		global.Log.Error("failed to articleTags : ", zap.Error(err))
		response.FailWithMessage("failed to articleTags", c)
	}
	response.OkWithData(tags, c)

}

// 文章收藏
func (articleApi *ArticleApi) ArticleLike(c *gin.Context) {
	var req request.ArticleLike
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserID = utils.GetUserID(c)
	err = articleService.ArticleLike(req)
	if err != nil {
		global.Log.Error("failed to articleLike : ", zap.Error(err))
		response.FailWithMessage("failed to articleLike", c)
	}
	response.OkWithMessage("Successfully articleLike", c)
}

// 判断是否文章收藏
func (articleApi *ArticleApi) ArticleIsLike(c *gin.Context) {
	var req request.ArticleLike
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserID = utils.GetUserID(c)
	isLike, err := articleService.ArticleIsLike(req)
	if err != nil {
		global.Log.Error("failed to articleIsLike : ", zap.Error(err))
		response.FailWithMessage("failed to articleIsLike", c)
	}
	response.OkWithData(isLike, c)

}

// 获取文章收藏列表
func (articleApi *ArticleApi) ArticleLikesList(c *gin.Context) {
	var pageInfo request.ArticleLikesList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	pageInfo.UserID = utils.GetUserID(c)
	list, total, err := articleService.ArticleLikesList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to articleLikesList :", zap.Error(err))
		response.FailWithMessage("Failed to articleLikesList", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// 创建文章
func (articleApi *ArticleApi) ArticleCreate(c *gin.Context) {
	var req request.ArticleCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = articleService.ArticleCreate(req)
	if err != nil {
		global.Log.Error("failed to article create:", zap.Error(err))
		response.FailWithMessage("failed to article create", c)
		return
	}

	response.OkWithMessage("Successfully articleCreate", c)
}

// 删除文章
func (articleApi *ArticleApi) ArticleDelete(c *gin.Context) {
	var req request.ArticleDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = articleService.ArticleDelete(req)
	if err != nil {
		global.Log.Error("failed to article delete:", zap.Error(err))
		response.FailWithMessage("failed to article delete", c)
		return
	}

	response.OkWithMessage("Successfully articleDelete", c)
}

// 更新文章
func (articleApi *ArticleApi) ArticleUpdate(c *gin.Context) {
	var req request.ArticleUpdate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = articleService.ArticleUpdate(req)
	if err != nil {
		global.Log.Error("failed to article update:", zap.Error(err))
		response.FailWithMessage("failed to article update", c)
		return
	}

	response.OkWithMessage("Successfully articleUpdate", c)

}

// 获取文章列表
func (articleApi *ArticleApi) ArticleList(c *gin.Context) {
	var pageInfo request.ArticleList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	articleList, total, err := articleService.ArticleList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to articleList", zap.Error(err))
		response.FailWithMessage("Failed to articleList", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  articleList,
		Total: total,
	}, c)
}

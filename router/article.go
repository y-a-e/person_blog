package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type ArticleRouter struct {
}

func (a *ArticleRouter) InitArticleRouter(publicRouter *gin.RouterGroup, privateRouter *gin.RouterGroup, adminRouter *gin.RouterGroup) {
	articleRouter := privateRouter.Group("article")
	articlePublicRouter := publicRouter.Group("article")
	articleAdminRouter := adminRouter.Group("article")

	articleApi := api.ApiGroupApp.ArticleApi
	{
		articleRouter.POST("like", articleApi.ArticleLike)          // 收藏
		articleRouter.GET("isLike", articleApi.ArticleIsLike)       // 是否收藏
		articleRouter.GET("likesList", articleApi.ArticleLikesList) // 收藏列表
	}
	{
		articlePublicRouter.GET(":id", articleApi.ArticleInfoByID)      // 根据id获取文章信息
		articlePublicRouter.GET("search", articleApi.ArticleSearch)     // 文章搜索
		articlePublicRouter.GET("category", articleApi.ArticleCategory) // 文章分类
		articlePublicRouter.GET("tags", articleApi.ArticleTags)         // 文章标签
	}
	{
		articleAdminRouter.POST("create", articleApi.ArticleCreate)   // 创建文章
		articleAdminRouter.DELETE("delete", articleApi.ArticleDelete) // 删除文章
		articleAdminRouter.PUT("update", articleApi.ArticleUpdate)    // 更新文章
		articleAdminRouter.GET("list", articleApi.ArticleList)        // 文章列表
	}
}

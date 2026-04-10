package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type CommentRouter struct {
}

func (c *CommentRouter) InitCommentRouter(publicRouter *gin.RouterGroup, privateRouter *gin.RouterGroup, adminRouter *gin.RouterGroup) {
	commentRouter := privateRouter.Group("comment")
	commentPublicRouter := publicRouter.Group("comment")
	commentAdminRouter := adminRouter.Group("comment")

	commentApi := api.ApiGroupApp.CommentApi
	{
		commentRouter.POST("create", commentApi.CommentCreate)   // 新增评论
		commentRouter.DELETE("delete", commentApi.CommentDelete) // 删除评论
		commentRouter.GET("info", commentApi.CommentInfo)        // 获取个人全部评论
	}
	{
		commentPublicRouter.GET(":article_id", commentApi.CommentInfoByArticleID) // 根据文章id获取评论内容
		commentPublicRouter.GET("new", commentApi.CommentNew)                     // 获取最新评论               // 文章搜索
	}
	{
		commentAdminRouter.GET("list", commentApi.CommentList) // 评论列表
	}
}

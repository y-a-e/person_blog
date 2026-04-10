package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type FeedbackRouter struct {
}

func (c *CommentRouter) InitFeedbackRouter(publicRouter *gin.RouterGroup, privateRouter *gin.RouterGroup, adminRouter *gin.RouterGroup) {
	feedbackRouter := privateRouter.Group("feedback")
	feedbackPublicRouter := publicRouter.Group("feedback")
	feedbackAdminRouter := adminRouter.Group("feedback")

	feedbackApi := api.ApiGroupApp.FeedbackApi
	{
		feedbackRouter.POST("create", feedbackApi.FeedbackCreate) // 新增反馈
		feedbackRouter.GET("info", feedbackApi.FeedbackInfo)      // 获取个人全部反馈
	}
	{
		feedbackPublicRouter.GET("new", feedbackApi.FeedbackNew) // 获取最新5条评论
	}
	{
		feedbackAdminRouter.DELETE("delete", feedbackApi.FeedbackDelete) // 删除反馈
		feedbackAdminRouter.PUT("reply", feedbackApi.FeedbackReply)      // 反馈返回
		feedbackAdminRouter.GET("list", feedbackApi.FeedbackList)        // 所有反馈列表
	}
}

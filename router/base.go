package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct {
}

func (b *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	baseApi := api.ApiGroupApp.BaseApi
	{
		baseRouter.POST("captcha", baseApi.Captcha)                                     // 图形验证码
		baseRouter.POST("sendEmailVerificationCode", baseApi.SendEmailVerificationCode) // 邮箱验证码
		baseRouter.GET("qqLoginURL", baseApi.QQLoginURL)                                // qq登录链接
	}
}

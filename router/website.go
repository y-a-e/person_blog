package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type WebsiteRouter struct {
}

func (w *WebsiteRouter) InitWebsiteRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	websiteRouter := Router.Group("website")
	websitePublicRouter := PublicRouter.Group("website")

	websiteApi := api.ApiGroupApp.WebsiteApi
	{
		websiteRouter.POST("addCarousel", websiteApi.WebsiteAddCarousel)             //	添加首页封面
		websiteRouter.PUT("cancelCarousel", websiteApi.WebsiteCancelCarousel)        // 取消首页封面
		websiteRouter.POST("createFooterLink", websiteApi.WebsiteCreateFooterLink)   // 创建页脚链接
		websiteRouter.DELETE("deleteFooterLink", websiteApi.WebsiteDeleteFooterLink) // 删除页脚链接
	}
	{
		websitePublicRouter.GET("logo", websiteApi.WebsiteLogo)             // 返回网站页签的图标
		websitePublicRouter.GET("title", websiteApi.WebsiteTitle)           // 返回网站页签的标题
		websitePublicRouter.GET("info", websiteApi.WebsiteInfo)             // 返回yaml-website信息
		websitePublicRouter.GET("carousel", websiteApi.WebsiteCarousel)     // 首页走马灯图片
		websitePublicRouter.GET("news", websiteApi.WebsiteNews)             // 获取每日新闻
		websitePublicRouter.GET("calendar", websiteApi.WebsiteCalendar)     // 获取日历信息
		websitePublicRouter.GET("footerLink", websiteApi.WebsiteFooterLink) // 获取页脚链接
	}
}

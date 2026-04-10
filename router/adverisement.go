package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type AdvertisementRouter struct {
}

func (a *AdvertisementRouter) InitAdvertisementRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	advertisementRouter := Router.Group("advertisement")
	advertisementPublicRouter := PublicRouter.Group("advertisement")

	adverisementApi := api.ApiGroupApp.AdverisementApi
	{
		advertisementRouter.POST("create", adverisementApi.AdvertisementCreate)   // 创建广告
		advertisementRouter.DELETE("delete", adverisementApi.AdvertisementDelete) // 删除广告
		advertisementRouter.PUT("update", adverisementApi.AdvertisementUpdate)    // 更新广告
		advertisementRouter.GET("list", adverisementApi.AdvertisementList)        // 广告列表
	}
	{
		advertisementPublicRouter.GET("info", adverisementApi.AdvertisementInfo) // 广告信息
	}
}

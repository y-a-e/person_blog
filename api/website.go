package api

import (
	"net/http"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WebsiteApi struct {
}

func (websiteApi *WebsiteApi) WebsiteAddCarousel(c *gin.Context) {
	var req request.WebsiteCarouselOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	err = websiteService.WebsiteAddCarousel(req)
	if err != nil {
		global.Log.Error("failed to websiteAddCarousel:", zap.Error(err))
		response.FailWithMessage("failed to websiteAddCarousel ", c)
	}
	response.OkWithMessage("Successfully websiteAddCarousel", c)
}

func (websiteApi *WebsiteApi) WebsiteCancelCarousel(c *gin.Context) {
	var req request.WebsiteCarouselOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	err = websiteService.WebsiteCancelCarousel(req)
	if err != nil {
		global.Log.Error("failed to websiteAddCarousel:", zap.Error(err))
		response.FailWithMessage("failed to websiteAddCarousel ", c)
	}
	response.OkWithMessage("Successfully websiteAddCarousel", c)
}

func (websiteApi *WebsiteApi) WebsiteLogo(c *gin.Context) {
	if global.Config.Website.Logo != "" {
		// c.Redirect,Gin框架的重定向方法;http.StatusMovedPermanently状态码301 永久重定向
		c.Redirect(http.StatusMovedPermanently, global.Config.Website.Logo)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/image/logo.png")
	}
}

func (websiteApi *WebsiteApi) WebsiteTitle(c *gin.Context) {
	// 封装网站标题 map[string]interface{}，参考response.go
	c.JSON(http.StatusOK, gin.H{
		"title": global.Config.Website.Title,
	})
}

func (websiteApi *WebsiteApi) WebsiteInfo(c *gin.Context) {
	// 直接读取配置返回网站信息
	response.OkWithData(global.Config.Website, c)
}

func (websiteApi *WebsiteApi) WebsiteCarousel(c *gin.Context) {
	urls, err := websiteService.WebsiteCarousel()
	if err != nil {
		global.Log.Error("failed to websiteCarousel:", zap.Error(err))
		response.FailWithMessage("failed to websiteCarousel ", c)
	}
	response.OkWithData(urls, c)
}

func (websiteApi *WebsiteApi) WebsiteNews(c *gin.Context) {
	sourceStr := c.Query("source")
	hotSearchData, err := websiteService.WebsiteNews(sourceStr)
	if err != nil {
		global.Log.Error("failed to websiteNews:", zap.Error(err))
		response.FailWithMessage("failed to websiteNews ", c)
	}
	response.OkWithData(hotSearchData, c)
}

func (websiteApi *WebsiteApi) WebsiteCalendar(c *gin.Context) {
	dateStr := time.Now().Format("2006/0102")
	calendar, err := websiteService.WebsiteCalendar(dateStr)
	if err != nil {
		global.Log.Error("failed to websiteCalendar:", zap.Error(err))
		response.FailWithMessage("failed to websiteCalendar ", c)
	}
	response.OkWithData(calendar, c)
}

func (websiteApi *WebsiteApi) WebsiteFooterLink(c *gin.Context) {
	footerLinkList := websiteService.WebsiteFooterLink()
	response.OkWithData(footerLinkList, c)
}

func (websiteApi *WebsiteApi) WebsiteCreateFooterLink(c *gin.Context) {
	var req database.FooterLink
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	err = websiteService.WebsiteCreateFooterLink(req)
	if err != nil {
		global.Log.Error("failed to websiteCreateFooterLink:", zap.Error(err))
		response.FailWithMessage("failed to websiteCreateFooterLink ", c)
	}
	response.OkWithMessage("Successfully websiteCreateFooterLink", c)
}

func (websiteApi *WebsiteApi) WebsiteDeleteFooterLink(c *gin.Context) {
	var req database.FooterLink
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	err = websiteService.WebsiteDeleteFooterLink(req)
	if err != nil {
		global.Log.Error("failed to wbsiteDeleteFooterLink:", zap.Error(err))
		response.FailWithMessage("failed to wbsiteDeleteFooterLink ", c)
	}
	response.OkWithMessage("Successfully wbsiteDeleteFooterLink", c)
}

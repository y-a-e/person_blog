package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdverisementApi struct {
}

func (adverisementApi *AdverisementApi) AdvertisementCreate(c *gin.Context) {
	var req request.AdvertisementCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	err = advertisementService.AdvertisementCreate(req)
	if err != nil {
		response.FailWithMessage("failed to advertisementCreate ", c)
	}
	response.OkWithMessage("Successfully advertisementCreate", c)
}

func (adverisementApi *AdverisementApi) AdvertisementDelete(c *gin.Context) {
	var req request.AdvertisementDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	err = advertisementService.AdvertisementDelete(req)
	if err != nil {
		response.FailWithMessage("failed to advertisementDelete ", c)
	}
	response.OkWithMessage("Successfully advertisementDelete", c)
}

func (adverisementApi *AdverisementApi) AdvertisementUpdate(c *gin.Context) {
	var req request.AdvertisementUpdate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = advertisementService.AdvertisementUpdate(req)
	if err != nil {
		global.Log.Error("failed to advertisement update:", zap.Error(err))
		response.FailWithMessage("failed to advertisement update", c)
		return
	}

	response.OkWithMessage("Successfully advertisementUpdate", c)
}

func (adverisementApi *AdverisementApi) AdvertisementList(c *gin.Context) {
	var req request.AdvertisementList
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	list, total, err := advertisementService.AdvertisementList(req)
	if err != nil {
		response.FailWithMessage("failed to advertisementList ", c)
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

func (adverisementApi *AdverisementApi) AdvertisementInfo(c *gin.Context) {
	list, total, err := advertisementService.AdvertisementInfo()
	if err != nil {
		global.Log.Error("failed to advertisementInfo ", zap.Error(err))
		response.FailWithMessage("failed to advertisementInfo ", c)
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

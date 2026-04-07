package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageApi struct {
}

func (imageApi *ImageApi) ImageUpload(c *gin.Context) {
	_, header, err := c.Request.FormFile("image")
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	url, err := imageService.ImageUpload(header)
	if err != nil {
		global.Log.Error("Failed to ImageUpload", zap.Error(err))
		response.FailWithMessage("Failed to ImageUpload", c)
		return
	}
	response.OkWithDetailed(response.ImageUpload{
		Url:     url,
		OssType: global.Config.System.OssType,
	}, "Successfully uploaded image", c)
}

func (imageApi *ImageApi) ImageDelete(c *gin.Context) {
	var req request.ImageDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = imageService.ImageDelete(req)
	if err != nil {
		global.Log.Error("Failed to ImageDelete", zap.Error(err))
		response.FailWithMessage("Failed to ImageDelete", c)
		return
	}

	response.OkWithMessage("Successfully uploaded image", c)
}

func (imageApi *ImageApi) ImageList(c *gin.Context) {
	var pageInfo request.ImageList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	imageList, total, err := imageService.ImageList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to ImageList", zap.Error(err))
		response.FailWithMessage("Failed to ImageList", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  imageList,
		Total: total,
	}, c)
}

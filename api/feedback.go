package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FeedbackApi struct {
}

func (feedbackApi *FeedbackApi) FeedbackCreate(c *gin.Context) {
	var req request.FeedbackCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	req.UUID = utils.GetUUID(c)
	err = feedbackService.FeedbackCreate(req)
	if err != nil {
		global.Log.Error("failed to feedbackCreate:", zap.Error(err))
		response.FailWithMessage("failed to feedbackCreate ", c)
	}
	response.OkWithMessage("Successfully feedbackCreate", c)
}

func (feedbackApi *FeedbackApi) FeedbackInfo(c *gin.Context) {
	uuid := utils.GetUUID(c)
	list, err := feedbackService.FeedbackInfo(uuid)
	if err != nil {
		global.Log.Error("Failed to feedbackInfo:", zap.Error(err))
		response.FailWithMessage("Failed to feedbackInfo", c)
		return
	}
	response.OkWithData(list, c)
}

func (feedbackApi *FeedbackApi) FeedbackNew(c *gin.Context) {
	list, err := feedbackService.FeedbackNew()
	if err != nil {
		global.Log.Error("Failed to get new feedback:", zap.Error(err))
		response.FailWithMessage("Failed to get new feedback", c)
		return
	}
	response.OkWithData(list, c)
}

func (feedbackApi *FeedbackApi) FeedbackDelete(c *gin.Context) {
	var req request.FeedbackDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	err = feedbackService.FeedbackDelete(req)
	if err != nil {
		global.Log.Error("Failed to feedbackDelete:", zap.Error(err))
		response.FailWithMessage("failed to feedbackDelete ", c)
	}
	response.OkWithMessage("Successfully feedbackDelete", c)
}

func (feedbackApi *FeedbackApi) FeedbackReply(c *gin.Context) {
	var req request.FeedbackReply
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	err = feedbackService.FeedbackReply(req)
	if err != nil {
		global.Log.Error("Failed to feedbackReply:", zap.Error(err))
		response.FailWithMessage("failed to feedbackReply ", c)
	}
	response.OkWithMessage("Successfully feedbackReply", c)
}

func (feedbackApi *FeedbackApi) FeedbackList(c *gin.Context) {
	var req request.PageInfo
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	list, total, err := feedbackService.FeedbackList(req)
	if err != nil {
		global.Log.Error("Failed to feedbackList:", zap.Error(err))
		response.FailWithMessage("failed to feedbackList ", c)
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

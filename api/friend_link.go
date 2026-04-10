package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FriendLinkApi struct {
}

func (friendLinkApi *FriendLinkApi) FriendLinkCreate(c *gin.Context) {
	var req request.FriendLinkCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	err = friendLinkService.FriendLinkCreate(req)
	if err != nil {
		response.FailWithMessage("failed to friendLinkCreate ", c)
	}
	response.OkWithMessage("Successfully friendLinkCreate", c)
}

func (friendLinkApi *FriendLinkApi) FriendLinkDelete(c *gin.Context) {
	var req request.FriendLinkDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	err = friendLinkService.FriendLinkDelete(req)
	if err != nil {
		response.FailWithMessage("failed to friendLinkDelete ", c)
	}
	response.OkWithMessage("Successfully friendLinkDelete", c)

}

func (friendLinkApi *FriendLinkApi) FriendLinkUpdate(c *gin.Context) {
	var req request.FriendLinkUpdate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = friendLinkService.FriendLinkUpdate(req)
	if err != nil {
		global.Log.Error("failed to friendLink update:", zap.Error(err))
		response.FailWithMessage("failed to friendLink update", c)
		return
	}

	response.OkWithMessage("Successfully friendLinkUpdate", c)
}

func (friendLinkApi *FriendLinkApi) FriendLinkList(c *gin.Context) {
	var req request.FriendLinkList
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	list, total, err := friendLinkService.FriendLinkList(req)
	if err != nil {
		response.FailWithMessage("failed to friendLinkList ", c)
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

func (friendLinkApi *FriendLinkApi) FriendLinkInfo(c *gin.Context) {
	list, total, err := friendLinkService.FriendLinkInfo()
	if err != nil {
		global.Log.Error("failed to friendLinkInfo ", zap.Error(err))
		response.FailWithMessage("failed to friendLinkInfo ", c)
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

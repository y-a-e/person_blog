package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentApi struct {
}

func (commentApi *CommentApi) CommentCreate(c *gin.Context) {
	var req request.CommentCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	req.UserUUID = utils.GetUUID(c)
	err = commentService.CommentCreate(req)
	if err != nil {
		response.FailWithMessage("failed to commentCreate ", c)
	}
	response.OkWithMessage("Successfully commentCreate", c)
}

func (commentApi *CommentApi) CommentDelete(c *gin.Context) {
	var req request.CommentDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	err = commentService.CommentDelete(c, req)
	if err != nil {
		response.FailWithMessage("failed to commentDelete ", c)
	}
	response.OkWithMessage("Successfully commentDelete", c)
}

func (commentApi *CommentApi) CommentInfo(c *gin.Context) {
	var uuid = utils.GetUUID(c)

	list, err := commentService.CommentInfo(uuid)
	if err != nil {
		response.FailWithMessage("failed to commentInfoByArticleId ", c)
	}
	response.OkWithData(list, c)
}

func (commentApi *CommentApi) CommentInfoByArticleID(c *gin.Context) {
	var req request.CommentInfoByArticleID
	err := c.ShouldBindUri(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	commentInfo, err := commentService.CommentInfoByArticleID(req)
	if err != nil {
		response.FailWithMessage("failed to commentInfoByArticleId ", c)
	}
	response.OkWithData(commentInfo, c)
}

func (commentApi *CommentApi) CommentNew(c *gin.Context) {
	list, err := commentService.CommentNew()
	if err != nil {
		response.FailWithMessage("failed to commentNew ", c)
	}
	response.OkWithData(list, c)
}

func (commentApi *CommentApi) CommentList(c *gin.Context) {
	var req request.CommentList
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

	list, total, err := commentService.CommentList(req)
	if err != nil {
		response.FailWithMessage("failed to commentList ", c)
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)

}

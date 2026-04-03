package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type BaseApi struct {
}

var store = base64Captcha.DefaultMemStore

func (baseApi BaseApi) Captcha(c *gin.Context) {
	// 调用github.com/mojocn/base64Captcha，创建图形验证码驱动
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)

	captcha := base64Captcha.NewCaptcha(driver, store)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		global.Log.Error("Failed to Captcha:", zap.Error(err))
		response.FailWithMessage("failed to Captcha", c)
		return
	}
	response.OkWithData(response.Captcha{
		CaptchaId: id,
		PicPath:   b64s,
	}, c)
}

func (baseApi BaseApi) SendEmailVerificationCode(c *gin.Context) {
	var req request.SendEmailVerificationCode
	err := c.ShouldBindJSON(&req) // 绑定对象
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证验证码答案
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		err := baseService.SendEmailVerificationCode(c, req.Email)
		if err != nil {
			global.Log.Error("Failed to SendEmail:", zap.Error(err))
			response.FailWithMessage("Failed to SendEmail", c)
			return
		}
		response.OkWithMessage("Successfully SendEmail", c)
		return
	}
	response.FailWithMessage("Failed to Email Verify err", c)
}

// QQLoginURL 返回 QQ 登录链接
func (baseApi *BaseApi) QQLoginURL(c *gin.Context) {
	url := global.Config.QQ.QQLoginURL()
	response.OkWithData(url, c)
}

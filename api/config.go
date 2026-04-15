package api

import (
	"server/config"
	"server/global"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ConfigApi struct {
}

func (configApi *ConfigApi) GetWebsite(c *gin.Context) {
	response.OkWithData(global.Config.Website, c)
}

func (configApi *ConfigApi) UpdateWebsite(c *gin.Context) {
	var req config.Website
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	err = configService.UpdateWebsite(req)
	if err != nil {
		response.FailWithMessage("failed to updateWebsite ", c)
	}
	response.OkWithMessage("Successfully updateWebsite", c)
}

func (configApi *ConfigApi) GetSystem(c *gin.Context) {
	response.OkWithData(global.Config.System, c)
}

func (configApi *ConfigApi) UpdateSystem(c *gin.Context) {
	var req config.System
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	err = configService.UpdateSystem(req)
	if err != nil {
		response.FailWithMessage("failed to updateSystem ", c)
	}
	response.OkWithMessage("Successfully updateSystem", c)
}

func (configApi *ConfigApi) GetEmail(c *gin.Context) {
	response.OkWithData(global.Config.Email, c)
}

func (configApi *ConfigApi) UpdateEmail(c *gin.Context) {
	var req config.Email
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	err = configService.UpdateEmail(req)
	if err != nil {
		response.FailWithMessage("failed to updateEmail ", c)
	}
	response.OkWithMessage("Successfully updateEmail", c)
}

func (configApi *ConfigApi) GetQQ(c *gin.Context) {
	response.OkWithData(global.Config.QQ, c)
}

func (configApi *ConfigApi) UpdateQQ(c *gin.Context) {
	var req config.QQ
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	err = configService.UpdateQQ(req)
	if err != nil {
		response.FailWithMessage("failed to updateQQ ", c)
	}
	response.OkWithMessage("Successfully updateQQ", c)
}

func (configApi *ConfigApi) GetQiniu(c *gin.Context) {
	response.OkWithData(global.Config.Qiniu, c)
}

func (configApi *ConfigApi) UpdateQiniu(c *gin.Context) {
	var req config.Qiniu
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	err = configService.UpdateQiniu(req)
	if err != nil {
		response.FailWithMessage("failed to updateQiniu ", c)
	}
	response.OkWithMessage("Successfully updateQiniu", c)
}

func (configApi *ConfigApi) GetJwt(c *gin.Context) {
	response.OkWithData(global.Config.Jwt, c)
}

func (configApi *ConfigApi) UpdateJwt(c *gin.Context) {
	var req config.Jwt
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	err = configService.UpdateJwt(req)
	if err != nil {
		response.FailWithMessage("failed to updateJwt ", c)
	}
	response.OkWithMessage("Successfully updateJwt", c)
}

func (configApi *ConfigApi) GetGaode(c *gin.Context) {
	response.OkWithData(global.Config.Gaode, c)
}

func (configApi *ConfigApi) UpdateGaode(c *gin.Context) {
	var req config.Gaode
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	err = configService.UpdateGaode(req)
	if err != nil {
		response.FailWithMessage("failed to updateGaode ", c)
	}
	response.OkWithMessage("Successfully updateGaode", c)
}

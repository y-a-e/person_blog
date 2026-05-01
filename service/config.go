package service

import (
	"server/config"
	"server/global"
	"server/model/appTypes"
	"server/utils"

	"gorm.io/gorm"
)

type ConfigService struct {
}

func (configService *ConfigService) UpdateWebsite(req config.Website) error {
	oldWebsite := []string{
		global.Config.Website.Logo,
		global.Config.Website.FullLogo,
		global.Config.Website.QQImage,
		global.Config.Website.WechatImage,
	}

	newWebsite := []string{
		req.Logo,
		req.FullLogo,
		req.QQImage,
		req.WechatImage,
	}

	addWebsite, removeWebsite := utils.DiffArrays(oldWebsite, newWebsite)
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.InitImagesCategory(tx, removeWebsite); err != nil {
			return err
		}

		if err := utils.ChangeImagesCategory(tx, addWebsite, appTypes.System); err != nil {
			return err
		}

		global.Config.Website = req
		if err := utils.SaveYAML(); err != nil {
			return err
		}

		return nil
	})
}

func (configService *ConfigService) UpdateSystem(req config.System) error {
	global.Config.System.UseMultipoint = req.UseMultipoint   // 是否启用多点登录拦截，防止同一账户在多个地方同时登录
	global.Config.System.SessionsSecret = req.SessionsSecret // 用于加密会话的密钥，确保会话数据的安全性
	global.Config.System.OssType = req.OssType               // 对应的对象存储服务类型，如 "local" 或 "qiniu"
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateEmail(req config.Email) error {
	global.Config.Email = req
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateQQ(req config.QQ) error {
	global.Config.QQ = req
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateQiniu(req config.Qiniu) error {
	global.Config.Qiniu = req
	return utils.SaveYAML()

}

func (configService *ConfigService) UpdateJwt(req config.Jwt) error {
	global.Config.Jwt = req
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateGaode(req config.Gaode) error {
	global.Config.Gaode = req
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateOss(req config.Oss) error {
	global.Config.Oss = req
	return utils.SaveYAML()
}

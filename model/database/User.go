package database

import (
	"server/global"
	"server/model/appTypes"

	"github.com/gofrs/uuid"
)

// 用户信息表
type User struct {
	global.MODEL
	UUID      uuid.UUID         `json:"uuid" gorm:"type:char(36);unique"`        // uuid
	Username  string            `json:"username"`                                // 用户名
	Password  string            `json:"-"`                                       // 传给前端时忽略 // 密码
	Email     string            `json:"email"`                                   // 邮箱
	OpenId    string            `json:"openId"`                                  // openid
	Avatar    string            `json:"avatar" gorm:"size:255"`                  // 头像
	Address   string            `json:"address"`                                 // 地址
	Signature string            `json:"signature" gorm:"default:'签名时空白，用户有点低调'"` // 签名
	RoleID    appTypes.RoleID   `json:"role_id"`                                 // 用户id
	Register  appTypes.Register `json:"register"`                                // 注册来源
	Freeze    bool              `json:"freeze"`                                  // 是否被冻结
}

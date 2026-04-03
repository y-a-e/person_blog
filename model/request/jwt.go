package request

import (
	"server/model/appTypes"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

// 自定义JWT格式，其中包括标准化数据
type JwtCustomClaims struct {
	BaseClaims
	jwt.RegisteredClaims // 标准JWT格式，包括过期时间、发行者
}

// 用户的基础信息
type BaseClaims struct {
	UserID uint            // 用户的ID
	UUID   uuid.UUID       // 用户的UUID
	RoleID appTypes.RoleID // 用户的级别
}

// 刷新JWT
type JwtCustomRefreshClaims struct {
	UserID               uint // 用户ID，用于与刷新Token相关的身份验证
	jwt.RegisteredClaims      // 标准JWT格式，包括过期时间、发行者
}

package database

import "server/global"

// jwt令牌的黑名单
type JwtBlacklist struct {
	global.MODEL        // 附加通用的数据库结构
	Jwt          string `json:"jwt" grom:"type:text"` // Jwt
}

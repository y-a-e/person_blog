package middleware

import (
	"server/model/appTypes"
	"server/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// 检查用户是否具有管理员权限
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := utils.GetRoleID(c)

		// 检查用户是否为管理员
		if roleID != appTypes.Admin {
			response.Forbidden("Access denied. Admin privileges are required", c)
			c.Abort()
			return
		}

		c.Next()
	}
}

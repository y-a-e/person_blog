package router

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (u *UserRouter) InitUserRouter(publicRouter *gin.RouterGroup, privateRouter *gin.RouterGroup, adminRouter *gin.RouterGroup) {
	userRouter := privateRouter.Group("user")
	userPublicRouter := publicRouter.Group("user")
	userLoginRouter := publicRouter.Group("user").Use(middleware.LoginRecord())
	userAdminRouter := adminRouter.Group("user")
	userApi := api.ApiGroupApp.UserApi
	{
		userRouter.POST("logout", userApi.Logout)                  // 登出
		userRouter.PUT("resetPassword", userApi.UserResetPassword) // 重设密码
		userRouter.GET("info", userApi.UserInfo)                   // 用户信息
		userRouter.PUT("changeInfo", userApi.UserChangeInfo)       // 修改用户信息
		userRouter.GET("weather", userApi.UserWeather)             // 用户天气
		userRouter.GET("chart", userApi.UserChart)                 // 用户图标数据
	}
	{
		userPublicRouter.POST("forgotPassword", userApi.ForgotPassword) // 忘记密码
		userPublicRouter.GET("card", userApi.UserCard)                  // 用户信息
	}
	{
		userLoginRouter.POST("register", userApi.Register) // 注册
		userLoginRouter.POST("login", userApi.Login)       // 登录
	}
	{
		userAdminRouter.GET("list", userApi.UserList)           // 用户列表
		userAdminRouter.PUT("freeze", userApi.UserFreeze)       // 冻结用户
		userAdminRouter.PUT("unfreeze", userApi.UserUnfreeze)   // 解冻用户
		userAdminRouter.GET("loginList", userApi.UserLoginList) // 登录日志列表
	}
}

package api

import (
	"errors"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type UserApi struct {
}

func (userApi *UserApi) Logout(c *gin.Context) {
	userService.Logout(c)
	response.OkWithMessage("Successful login", c)

}

func (userApi *UserApi) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserID = utils.GetUserID(c)

	err = userService.UserResetPassword(req)
	if err != nil {
		global.Log.Error("failed to reset password error:", zap.Error(err))
		response.FailWithMessage("failed to modify", c)
		return
	}
	response.OkWithMessage("Successfully changed password, please log in again", c)
	userService.Logout(c)
}

func (userApi *UserApi) UserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	userInto, err := userService.UserInfo(userID)
	if err != nil {
		global.Log.Error("failed to get into error:", zap.Error(err))
		response.FailWithMessage("failed to get", c)
		return
	}
	response.OkWithData(userInto, c)
}

func (userApi *UserApi) UserChangeInfo(c *gin.Context) {
	var req request.UserChangeInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserID = utils.GetUserID(c)
	err = userService.UserChangeInfo(req)
	if err != nil {
		global.Log.Error("Failed to change user information:", zap.Error(err))
		response.FailWithMessage("Failed to change user information", c)
		return
	}
	response.OkWithMessage("Successfully changed user information", c)
}

func (userApi *UserApi) UserWeather(c *gin.Context) {
	ip := c.ClientIP()
	weather, err := userService.UserWeather(ip)
	if err != nil {
		global.Log.Error("Failed to get user weather", zap.Error(err))
		response.FailWithMessage("Failed to get user weather", c)
		return
	}
	response.OkWithData(weather, c)
}

func (userApi *UserApi) UserChart(c *gin.Context) {
	var req request.UserChart
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := userService.UserChart(req)
	if err != nil {
		global.Log.Error("Failed to get user chart:", zap.Error(err))
		response.FailWithMessage("Failed to user chart", c)
		return
	}
	response.OkWithData(data, c)
}

func (userApi *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)

	// 验证邮箱
	email := session.Get("eamil")
	if email == nil || (email).(string) != req.Email {
		response.FailWithMessage("email error", c)
		return
	}

	// 验证验证码
	code := session.Get("verification_code")
	if code == nil || code.(string) != req.VerificationCode {
		response.FailWithMessage("code error", c)
		return
	}

	// 验证过期时间戳
	expireTime := session.Get("expire_time")
	if expireTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("expireTime error", c)
		return
	}

	err = userService.ForgotPassword(req)
	if err != nil {
		global.Log.Error("failed to register error:", zap.Error(err))
		response.FailWithMessage("register user err", c)
		return
	}
	response.OkWithMessage("Successfully retrieved", c)
}

func (userApi *UserApi) UserCard(c *gin.Context) {
	var req request.UserCard
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userCard, err := userService.UserCard(req)
	if err != nil {
		global.Log.Error("Failed to get card:", zap.Error(err))
		response.FailWithMessage("Failed to userCard error", c)
		return
	}

	response.OkWithData(userCard, c)
}

func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)

	// 验证邮箱
	email := session.Get("email")
	if email == nil || (email).(string) != req.Email {
		response.FailWithMessage("email error", c)
		return
	}

	// 验证验证码
	code := session.Get("verification_code")
	if code == nil || code.(string) != req.VerificationCode {
		response.FailWithMessage("code error", c)
		return
	}

	// 验证过期时间戳
	expireTime := session.Get("expire_time")
	if expireTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("expireTime error", c)
		return
	}

	// 注册
	u := database.User{Username: req.Username, Password: req.Password, Email: req.Email}
	user, err := userService.Register(u)
	if err != nil {
		global.Log.Error("failed to register error:", zap.Error(err))
		response.FailWithMessage("register user err", c)
		return
	}

	// 登录成功后生成 token
	userApi.TokenNext(c, user)
}

func (userApi *UserApi) Login(c *gin.Context) {
	switch c.Query("flag") {
	case "email":
		userApi.EmailLogin(c)
	case "qqlogin":
		userApi.QQLogin(c)
	default:
		userApi.EmailLogin(c)
	}
}

func (userApi *UserApi) EmailLogin(c *gin.Context) {
	var req request.Login
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 校验验证码
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		u := database.User{Email: req.Email, Password: req.Password}
		user, err := userService.EmailLogin(u)
		if err != nil {
			global.Log.Error("Failed to login:", zap.Error(err))
			response.FailWithMessage("Failed to login", c)
			return
		}

		// 登录成功后生成 token
		userApi.TokenNext(c, user)
		return
	}
	response.FailWithMessage("Incorrect verification code", c)
}

func (userApi *UserApi) QQLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("Code is required", c)
		return
	}

	// 获取访问令牌
	accessTokenResponse, err := qqService.GetAccessTokenByCode(code)
	if err != nil || accessTokenResponse.Openid == "" {
		global.Log.Error("Invalid code", zap.Error(err))
		response.FailWithMessage("Invalid code", c)
		return
	}

	// 根据访问令牌进行QQ登录
	user, err := userService.QQLogin(accessTokenResponse)
	if err != nil {
		global.Log.Error("Failed to login:", zap.Error(err))
		response.FailWithMessage("Failed to login", c)
		return
	}

	// 登录成功后生成 token
	userApi.TokenNext(c, user)
}

func (userApi *UserApi) TokenNext(c *gin.Context, user database.User) {
	// 检查用户是否被冻结
	if user.Freeze {
		response.FailWithMessage("The user is frozen, contact the administrator", c)
		return
	}

	baseClaims := request.BaseClaims{
		UserID: user.ID,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	}

	j := utils.NewJWT()

	// 创建访问令牌
	accessClaims := j.CreateAccessClaims(baseClaims)
	accessToken, err := j.CreateAccessToken(accessClaims)
	if err != nil {
		global.Log.Error("Failed to get accessToken:", zap.Error(err))
		response.FailWithMessage("Failed to get accessToken", c)
		return
	}

	// 创建刷新令牌
	refreshClaims := j.CreateRefreshClaims(baseClaims)
	refreshToken, err := j.CreateRefreshToken(refreshClaims)
	if err != nil {
		global.Log.Error("Failed to get refreshToken:", zap.Error(err))
		response.FailWithMessage("Failed to get refreshToken", c)
		return
	}

	// 是否开启了多地点登录拦截
	if !global.Config.System.UseMultipoint {
		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
		return
	}

	// 检查 Redis 中是否已存在该用户的 JWT
	if jwtStr, err := jwtService.GetRedisJWT(user.UUID); errors.Is(err, redis.Nil) {
		// 不存在就设置新的
		if err := jwtService.SetRedisJWT(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	} else if err != nil {
		// 出现错误处理
		global.Log.Error("Failed to set login status:", zap.Error(err))
		response.FailWithMessage("Failed to set login status", c)
	} else {
		// Redis 中已存在该用户的 JWT，将旧的 JWT 加入黑名单，并设置新的 token
		var blacklist database.JwtBlacklist
		blacklist.Jwt = jwtStr
		if err := jwtService.JoinInBlacklist(blacklist); err != nil {
			global.Log.Error("Failed to invalidate jwt:", zap.Error(err))
			response.FailWithMessage("Failed to invalidate jwt", c)
			return
		}

		// 设置新的 JWT 到 Redis
		if err := jwtService.SetRedisJWT(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	}
}

func (userApi *UserApi) UserList(c *gin.Context) {
	var pageInfo request.UserList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to getLists :", zap.Error(err))
		response.FailWithMessage("Failed to getLists", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

func (userApi *UserApi) UserFreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = userService.UserFreeze(req)
	if err != nil {
		global.Log.Error("Failed to user freeze :", zap.Error(err))
		response.FailWithMessage("Failed to freeze", c)
		return
	}
	response.OkWithMessage("Successfully freeze", c)
}

func (userApi *UserApi) UserUnfreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = userService.UserUnfreeze(req)
	if err != nil {
		global.Log.Error("Failed to user freeze :", zap.Error(err))
		response.FailWithMessage("Failed to unFreeze", c)
		return
	}
	response.OkWithMessage("Successfully unFreeze", c)
}

func (userApi *UserApi) UserLoginList(c *gin.Context) {
	var pageInfo request.UserLoginList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserLoginList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to getLists :", zap.Error(err))
		response.FailWithMessage("Failed to getLists", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

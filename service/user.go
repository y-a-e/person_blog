package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
}

func (userSrver *UserService) Logout(c *gin.Context) {
	uuid := utils.GetUUID(c)
	global.Redis.Del(uuid.String())
	jwt := utils.GetRefreshToken(c)
	utils.ClearRefreshToken(c)
	err := ServiceGroupApp.JoinInBlacklist(database.JwtBlacklist{Jwt: jwt})
	if err != nil {
		global.Log.Error("failed to register error:", zap.Error(err))
	}
}

func (userSrver *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}

	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("original password does not match the current account")
	}

	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userSrver *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (userSrver *UserService) UserChangeInfo(req request.UserChangeInfo) error {
	// var user database.User
	// if err := global.DB.Where("uuid = ?", req.UserID).First(&user).Error; err != nil {
	// 	return err
	// }

	// user.Username = req.Username
	// user.Address = req.Address
	// user.Signature = req.Signature
	// global.DB.Save(&user)

	// return nil

	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	return global.DB.Model(&user).Updates(req).Error
}

func (userSrver *UserService) UserWeather(ip string) (string, error) {
	// 从redis中获取天气数据，如果没有数据，则调用高德api进行查询
	result, err := global.Redis.Get("weather-" + ip).Result()
	if err != nil {
		ipResponse, err := ServiceGroupApp.GaodeService.GetLocationByIP(ip)
		if err != nil {
			return "", err
		}
		live, err := ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode)
		if err != nil {
			return "", err
		}

		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"

		// 将天气数据存入redis
		if err := global.Redis.Set("weather-"+ip, weather, time.Hour*1).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	return result, nil
}

func (userSrver *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	// 构建查询条件
	where := global.DB.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", req.Date))

	var res response.UserChart

	// 生成日期列表
	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i <= req.Date; i++ {
		res.DateList = append(res.DateList, startDate.AddDate(0, 0, i).Format("2006-03-01"))
	}
	// 获取登录数据
	loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
	// 获取注册数据
	registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)

	for _, date := range res.DateList {
		loginCount := loginCounts[date]
		registerCount := registerCounts[date]
		res.LoginData = append(res.LoginData, loginCount)
		res.RegisterData = append(res.RegisterData, registerCount)
	}

	return res, nil
}

func (userSrver *UserService) ForgotPassword(req request.ForgotPassword) error {
	var user database.User

	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil
	}
	user.Password = utils.BcryptHash(req.NewPassword)

	return global.DB.Save(user).Error
}

func (userSrver *UserService) UserCard(req request.UserCard) (response.UserCard, error) {
	var user database.User
	if err := global.DB.Where("uuid = ?", req.UUID).Select("uuid", "username", "avatar", "address", "signature").First(&user).Error; err != nil {
		return response.UserCard{}, err
	}

	return response.UserCard{
		UUID:      user.UUID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

func (userSrver *UserService) Register(u database.User) (database.User, error) {
	// 是否已注册
	if errors.Is(global.DB.Where("email = ?", u.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("email addresss is already")
	}

	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	u.Avatar = "/image/avatar.jpg"
	u.RoleID = appTypes.User
	u.Register = appTypes.Email

	// 数据库创建用户
	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	}

	return u, nil
}

func (userService *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return database.User{}, errors.New("incorrect email or password")
		}
		return user, nil
	}
	return database.User{}, err
}

func (userService *UserService) QQLogin(accessTokenResponse other.AccessTokenResponse) (database.User, error) {
	var user database.User

	// 尝试查找用户
	err := global.DB.Where("openid = ?", accessTokenResponse.Openid).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.User{}, err
	}

	// 如果用户不存在，则创建新用户
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userInfoResponse, err := ServiceGroupApp.QQService.GetUserInfoByAccessTokenAndOpenid(accessTokenResponse.AccessToken, accessTokenResponse.Openid)
		if err != nil {
			return database.User{}, err
		}
		user.UUID = uuid.Must(uuid.NewV4())
		user.Username = userInfoResponse.Nickname
		user.OpenId = accessTokenResponse.Openid
		user.Avatar = userInfoResponse.FigureurlQQ2
		user.RoleID = appTypes.User
		user.Register = appTypes.QQ

		if err := global.DB.Create(&user).Error; err != nil {
			return database.User{}, err
		}
	}

	return user, nil
}

func (userSrver *UserService) UserList(info request.UserList) (interface{}, int64, error) {
	db := global.DB

	if info.UUID != nil {
		db = db.Where("uuid = ?", info.UUID)
	}

	if info.Register != nil {
		db = db.Where("register = ?", info.Register)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.User{}, option)
}

func (userSrver *UserService) UserFreeze(req request.UserOperation) error {
	var user database.User
	// 冻结账号
	if err := global.DB.Take(&user, req.ID).Update("freeze", true).Error; err != nil {
		return err
	}

	// 利用jwt处理正处于登录状态的账号
	jwtStr, _ := ServiceGroupApp.JwtService.GetRedisJWT(user.UUID)
	if jwtStr != "" {
		_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
	}

	return nil
}

func (userSrver *UserService) UserUnfreeze(req request.UserOperation) error {
	return global.DB.Take(database.User{}, req.ID).Update("freeze", false).Error
}

func (userSrver *UserService) UserLoginList(info request.UserLoginList) (interface{}, int64, error) {
	db := global.DB

	if info.UUID != nil {
		var userID uint
		//	// update all users's name to `hello`
		//	db.Model(&User{}).Update("name", "hello")
		if err := global.DB.Model(&database.User{}).Update("id", userID).Error; err != nil {
			return nil, 0, err
		}

		db = db.Where("user_id = ?", userID)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
		Preload:  []string{"User"}, // Login表关联User表，查询时需要显示
	}

	return utils.MySQLPagination(&database.Login{}, option)
}

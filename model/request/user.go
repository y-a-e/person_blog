package request

type Register struct {
	Username         string `json:"username" binding:"required,max=20"`         // 用户名
	Password         string `json:"password" binding:"required,min=8,max=16"`   // 密码
	Email            string `json:"email" binding:"required,email"`             // 邮箱
	VerificationCode string `json:"verification_code" binding:"required,len=6"` // 验证码
}

type Login struct {
	Email     string `json:"email" binding:"required,email"`           // 邮箱
	Password  string `json:"password" binding:"required,min=8,max=16"` // 密码
	Captcha   string `json:"captcha" binding:"required,len=6"`         // 验证码
	CaptchaID string `json:"captcha_id" binding:"required"`            // 验证码ID
}

type ForgotPassword struct {
	Email            string `json:"email" binding:"required,email"`               // 邮箱
	VerificationCode string `json:"verification_code" binding:"required,len=6"`   // 验证码Password         string `json:"password" binding:"required,min=8,max=16"`   // 密码
	NewPassword      string `json:"new_password" binding:"required,min=8,max=16"` // 密码
}

type UserCard struct {
	UUID string `json:"uuid" form:"uuid" binding:"required"`
}

type UserResetPassword struct {
	UserID      uint   `json:"-"`
	Password    string `json:"password" binding:"required,min=8,max=16"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=16"`
}

type UserChangeInfo struct {
	UserID    uint   `json:"-"`
	Username  string `json:"username" binding:"required,max=20"`
	Address   string `json:"address" binding:"max=200"`
	Signature string `json:"signature" binding:"max=320"`
}

type UserChart struct {
	Date int `json:"date" form:"date" binding:"required,oneof=7 30 90 180 365"`
}

type UserList struct {
	UUID     *string `json:"uuid" form:"uuid"`         // 用户id
	Register *string `json:"register" form:"register"` // 注册来源
	PageInfo         // 页码/页脚
}

type UserOperation struct {
	ID uint `json:"id" binding:"required"`
}

type UserLoginList struct {
	UUID *string `json:"uuid" form:"uuid"`
	PageInfo
}

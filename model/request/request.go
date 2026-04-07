package request

// 邮箱验证码
type SendEmailVerificationCode struct {
	Email     string `json:"email" binding:"required,email"`   // required,前端必传，email为格式，自动检查
	Captcha   string `json:"captcha" binding:"required,len=6"` // required,前端必传，len=6，验证码
	CaptchaID string `json:"captcha_id" binding:"required"`    // required,前端必传，验证码ID
}

type ImageDelete struct {
	IDs []uint `json:"ids"`
}

type ImageList struct {
	Name     *string `json:"name" form:"name"`
	Category *string `json:"category" form:"category"`
	Storage  *string `json:"storage" form:"storage"`
	PageInfo
}

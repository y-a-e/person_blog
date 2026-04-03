package response

// 图形验证码
type Captcha struct {
	CaptchaId string `json:"captcha_id"`
	PicPath   string `json:"pic_path"`
}

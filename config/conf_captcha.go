package config

type Captcha struct {
	Height   int     `json:"height" yaml:"height"`       // 高度
	Width    int     `json:"width" yaml:"width"`         // 宽度
	Length   int     `json:"length" yaml:"length"`       // 数字个数
	MaxSkew  float64 `json:"max_skew" yaml:"max_skew"`   // 倾斜度
	DotCount int     `json:"dot_count" yaml:"dot_count"` // 背景圆点数量
}

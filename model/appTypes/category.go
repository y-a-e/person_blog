package appTypes

import (
	"encoding/json"
)

type Category int

const (
	Null         Category = iota // 未使用
	System                       // 系统
	Carousel                     // 背景
	Cover                        // 封面
	Illustration                 // 插图
	AdImage                      // 广告
	Logo                         // 友链
)

// MarshalJSON 实现了 json.Marshaler 接口，字符串转化为json
func (c Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口，json转化为字符串
func (c *Category) UmmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*c = ToCategory(str)
	return nil
}

// 字符串转换为Category
func ToCategory(str string) Category {
	switch str {
	case "未使用":
		return Null
	case "系统":
		return System
	case "背景":
		return Carousel
	case "封面":
		return Cover
	case "插图":
		return Illustration
	case "广告":
		return AdImage
	case "友链":
		return Logo
	default:
		return -1

	}
}

// Category对象转换为字符串
func (c Category) String() string {
	switch c {
	case Null:
		return "未使用"
	case System:
		return "系统"
	case Carousel:
		return "背景"
	case Cover:
		return "封面"
	case Illustration:
		return "插图"
	case AdImage:
		return "广告"
	case Logo:
		return "友链"
	default:
		return "未知类别"
	}
}

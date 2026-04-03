package config

type Gaode struct {
	Enable bool   `json:"enable" yaml:"enable"` // 是否启动高德地图
	Key    string `json:"key" yaml:"key"`       // 高德密钥
}

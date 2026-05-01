package config

// Oss 阿里云OSS配置，详情请见 https://help.aliyun.com/document_detail/31837.html
type Oss struct {
	AccessKey string `json:"access_key" yaml:"access_key"` // 阿里云 AccessKey ID
	SecretKey string `json:"secret_key" yaml:"secret_key"` // 阿里云 AccessKey Secret
	Bucket    string `json:"bucket" yaml:"bucket"`         // Bucket 名称
	Endpoint  string `json:"endpoint" yaml:"endpoint"`     // 公网 Endpoint（如 oss-cn-guangzhou.aliyuncs.com）
	Region    string `json:"region" yaml:"region"`         // 地域代码（如 cn-hangzhou）
	ImgPath   string `json:"img_path" yaml:"img_path"`     // CDN 加速域名或公网访问地址
	UseHTTPS  bool   `json:"use_https" yaml:"use_https"`   // 是否使用 HTTPS
	ChunkSize int    `json:"chunk_size" yaml:"chunk_size"` // 分片上传大小（MB）
}

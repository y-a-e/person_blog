package upload

import (
	"mime/multipart"
	"server/global"
	"server/model/appTypes"
)

// 定义一个白名单映射，包含支持的图片文件类型
var WhiteImageList = map[string]struct{}{
	".jpg":  {},
	".png":  {},
	".jpeg": {},
	".ico":  {},
	".tiff": {},
	".gif":  {},
	".svg":  {},
	".webp": {},
}

type OSS interface {
	UploadImage(file *multipart.FileHeader) (string, string, error)
	DeleteImage(key string) error
}

// 根据配置中的 OssType 来选择使用的存储类型
func NewOss() OSS {
	switch global.Config.System.OssType {
	case "local":
		return &Local{}
	case "qinqu":
		return &Qiniu{}
	case "oss":
		return &Oss{}
	default:
		return &Local{}
	}
}

// 根据传参返回相应的 OSS 实例
func NewOssWithStorage(storage appTypes.Storage) OSS {
	switch storage {
	case appTypes.Local:
		return &Local{}
	case appTypes.Qiniu:
		return &Qiniu{}
	case appTypes.Oss:
		return &Oss{}
	default:
		return &Local{}
	}
}

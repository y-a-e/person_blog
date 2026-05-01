package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"server/global"
	"server/utils"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

type Oss struct {
}

// UploadImage 上传图片到阿里云OSS
// 返回值: 完整访问URL, 文件存储Key, 错误信息
func (*Oss) UploadImage(file *multipart.FileHeader) (string, string, error) {
	// ========== 步骤1: 检查文件大小 ==========
	// 将文件大小从字节转换为MB
	size := float64(file.Size) / float64(1024*1024)
	// 判断是否超过配置文件中的大小限制
	if size >= float64(global.Config.Upload.Size) {
		return "", "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Config.Upload.Size)
	}

	// ========== 步骤2: 检查文件类型 ==========
	// 获取文件扩展名（如 .jpg, .png）
	ext := filepath.Ext(file.Filename)
	// 去除扩展名，得到原始文件名
	name := strings.TrimSuffix(file.Filename, ext)
	// 验证文件扩展名是否在白名单中（只允许上传图片类型）
	if _, exists := WhiteImageList[ext]; !exists {
		return "", "", errors.New("don't upload files that aren't image types")
	}

	// ========== 步骤3: 创建OSS客户端 ==========
	// 初始化阿里云OSS客户端，用于后续API调用。OSS没有Token，直接使用AccessKey和SecretKey认证
	client, err := newOssClient()
	if err != nil {
		return "", "", fmt.Errorf("failed to create OSS client: %w", err)
	}

	// ========== 步骤4: 生成唯一存储键(文件名) ==========
	// 使用MD5哈希原文件名 + 时间戳生成唯一名称，避免文件名冲突
	// 格式: MD5(原文件名) + "-" + 时间(YYYYMMDDHHMMSS) + 扩展名
	// 示例: a1b2c3d4e5f6-20260501120000.jpg
	fileKey := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext

	// ========== 步骤5: 读取文件数据 ==========
	// 打开上传的文件，获取文件句柄
	data, err := file.Open()
	if err != nil {
		return "", "", err
	}
	// 确保函数返回前关闭文件句柄，释放资源
	defer data.Close()

	// ========== 步骤6: 上传文件到OSS ==========
	// 调用OSS的PutObject接口，将文件内容上传到指定的Bucket中
	// Bucket: 存储空间名称
	// Key: 文件在OSS中的存储路径/名称（即fileKey）
	// Body: 文件数据流
	_, err = client.PutObject(context.Background(), &oss.PutObjectRequest{
		Bucket: oss.Ptr(global.Config.Oss.Bucket), // Bucket名称，从配置获取
		Key:    oss.Ptr(fileKey),                  // 文件存储键，唯一标识文件
		Body:   data,                              // 文件数据
	})
	if err != nil {
		// 上传失败，返回错误信息
		return "", "", fmt.Errorf("failed to upload file to OSS: %w", err)
	}

	// ========== 步骤7: 返回结果 ==========
	// 返回完整访问URL（如 https://cdn.example.com/a1b2c3d4e5f6-20260501120000.jpg）
	// 以及文件存储键（用于删除文件时指定）
	return global.Config.Oss.ImgPath + fileKey, fileKey, nil
}

// DeleteImage 从阿里云OSS删除指定文件
// 参数: key 文件存储键（即上传时返回的fileKey）
// 返回值: 错误信息
func (*Oss) DeleteImage(key string) error {
	// ========== 步骤1: 创建OSS客户端 ==========
	// 初始化阿里云OSS客户端
	client, err := newOssClient()
	if err != nil {
		return fmt.Errorf("failed to create OSS client: %w", err)
	}

	// ========== 步骤2: 删除OSS上的文件 ==========
	// 调用OSS的DeleteObject接口，删除指定Bucket中的文件
	// Bucket: 存储空间名称
	// Key: 要删除的文件存储键
	_, err = client.DeleteObject(context.Background(), &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(global.Config.Oss.Bucket), // Bucket名称
		Key:    oss.Ptr(key),                      // 要删除的文件键
	})
	if err != nil {
		// 删除失败，返回错误信息
		return fmt.Errorf("failed to delete file from OSS: %w", err)
	}

	// 删除成功，返回nil
	return nil
}

// newOssClient 创建并返回阿里云OSS客户端
// 使用AK/SK凭证初始化客户端，建立与OSS服务的连接
func newOssClient() (*oss.Client, error) {
	// ========== 步骤1: 处理Endpoint协议头 ==========
	// 确保Endpoint包含协议头（http:// 或 https://）
	endpoint := global.Config.Oss.Endpoint
	// 如果没有协议头，根据UseHTTPS配置添加
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		if global.Config.Oss.UseHTTPS {
			// 使用HTTPS安全协议
			endpoint = "https://" + endpoint
		} else {
			// 使用HTTP协议
			endpoint = "http://" + endpoint
		}
	}

	// ========== 步骤2: 确定地域代码 ==========
	// 映射友好名称到阿里云Region代码
	region := global.Config.Oss.Region
	switch region {
	case "cn-guangzhou", "Guangzhou":
		region = "cn-guangzhou"
	case "cn-shenzhen", "Shenzhen":
		region = "cn-shenzhen"
	case "cn-hangzhou", "Hangzhou":
		region = "cn-hangzhou"
	case "cn-shanghai", "Shanghai":
		region = "cn-shanghai"
	case "cn-beijing", "Beijing":
		region = "cn-beijing"
	case "cn-qingdao", "Qingdao":
		region = "cn-qingdao"
	case "cn-hongkong", "Hongkong":
		region = "cn-hongkong"
	case "ap-southeast-1", "Singapore":
		region = "ap-southeast-1"
	case "us-west-1", "USWest":
		region = "us-west-1"
	case "us-east-1", "USEast":
		region = "us-east-1"
	}

	// ========== 步骤3: 配置OSS客户端 ==========
	// LoadDefaultConfig: 加载默认配置
	// WithCredentialsProvider: 设置访问凭证（AK/SK）
	// WithEndpoint: 设置OSS访问地址（外网Endpoint）
	// WithRegion: 设置地域（用于签名）
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			global.Config.Oss.AccessKey, // AccessKey ID
			global.Config.Oss.SecretKey, // AccessKey Secret
			"",                          // Token（临时凭证，可为空）
		)).
		WithEndpoint(endpoint). // OSS服务地址
		WithRegion(region)      // 地域代码

	// ========== 步骤4: 创建并返回OSS客户端 ==========
	return oss.NewClient(cfg), nil
}

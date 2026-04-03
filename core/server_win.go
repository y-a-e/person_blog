//go:build windows
// +build windows

// ⬆头部为指明为windows系统才使用该初始化，因为执行的go get "github.com/fvbock/endless"，是针对非windows环境的优雅重启
package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// initServer 函数初始化一个标准的 HTTP 服务器（适用于 Windows 系统）
func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,          // 设置服务器监听的地址
		Handler:        router,           // 设置请求处理器（路由）
		ReadTimeout:    10 * time.Minute, // 设置请求的读取超时时间为 10 分钟
		WriteTimeout:   10 * time.Minute, // 设置响应的写入超时时间为 10 分钟
		MaxHeaderBytes: 1 << 20,          // 设置最大请求头的大小（1MB）
	}
}

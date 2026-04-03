package main

import (
	"server/core"
	"server/global"
	"server/initialize"
)

func main() {
	global.Config = core.InitConf()          // core目录下加载全局配置
	global.Log = core.InitLogger()           // core目录下加载日志配置
	initialize.OtherInit()                   // initialize目录下创建定时任务器
	global.ESClient = initialize.ConnectES() // initialize目录下启动ES功能
	global.Redis = initialize.ConnectRedis() // initialize目录下启动Redis功能
	global.DB = initialize.InitGorm()        // initialize目录下数据库连接

	core.RunServer() // 启动服务器
}

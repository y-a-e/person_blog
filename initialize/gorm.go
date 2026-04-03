package initialize

import (
	"os"
	"server/global"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm() *gorm.DB {
	// 数据库连接
	mysqlCfg := global.Config.Mysql
	db, err := gorm.Open(mysql.Open(mysqlCfg.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(mysqlCfg.LogLevel()), // 设置日志级别
	})
	if err != nil {
		global.Log.Error("failed to connect to Mysql:", zap.Error(err))
		os.Exit(1)
	}

	// 获取DB连接池对象，并设置最大空闲和最大打开的连接数
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)

	return db
}

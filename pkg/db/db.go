package db

import (
	"go-my-blog/config"
	"go-my-blog/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化 MySQL 连接（使用解析后的配置）
func Init() {
	// 从全局配置中获取 MySQL 配置
	mysqlConf := config.Conf.Mysql

	// 配置 GORM 日志模式（开发环境打印 SQL，生产环境仅打印错误）
	var gormLog gormLogger.Interface
	if mysqlConf.LogMode {
		gormLog = gormLogger.Default.LogMode(gormLogger.Info)
	} else {
		gormLog = gormLogger.Default.LogMode(gormLogger.Error)
	}

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(mysqlConf.DSN), &gorm.Config{
		Logger: gormLog,
	})
	if err != nil {
		logger.Fatal("MySQL 连接失败", zap.Error(err), zap.String("dsn", mysqlConf.DSN))
	}

	// 配置连接池（使用配置中的参数）
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(mysqlConf.MaxOpenConns)
	sqlDB.SetMaxIdleConns(mysqlConf.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(mysqlConf.GetConnMaxLifetime()) // 使用辅助方法转时间

	DB = db
	logger.Info("MySQL 连接初始化成功")
}

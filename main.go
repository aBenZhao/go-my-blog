package main

import (
	"fmt"
	"go-my-blog/config"
	priorityConfig "go-my-blog/config/priority_config"
	"go-my-blog/pkg/db"
	"go-my-blog/pkg/logger"
	"go-my-blog/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// 1. 先初始化日志（确保后续错误能正常打印）
	priorityConfig.PriorityConfInit()
	logger.Init()
	// 程序退出时同步日志缓冲区
	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Error("日志同步失败", zap.Error(err))
		}
	}()

	// 2. 初始化配置（必须在数据库初始化前调用，因为数据库需要用配置中的 DSN）
	logger.Info("开始初始化配置")
	// 手动调用配置初始化函数
	config.Init()

	// 3. 初始化数据库（依赖配置中的 MySQL 参数）
	logger.Info("开始初始化数据库")
	db.Init()

	// 4. 初始化 Gin 引擎和路由
	logger.Info("开始初始化路由")
	r := gin.Default()
	router.InitRouter(r)

	// 5. 启动服务（依赖配置中的端口参数）
	port := config.Conf.Server.Port
	logger.Info("服务启动成功", zap.Int("port", port))
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Fatal("服务启动失败", zap.Error(err))
	}
}

package bootstrap

import (
	"go-my-blog/internal/handler"
	"go-my-blog/internal/repo"
	"go-my-blog/internal/service"

	"gorm.io/gorm"
)

// InitUserModule 初始化用户模块，创建并配置用户处理器
// @param db 数据库连接对象
// @return *handler.UserHandler 用户处理器实例
func InitUserModule(db *gorm.DB) *handler.UserHandler {
	// 创建用户仓库实例，传入数据库连接
	repository := repo.NewUserRepository(db)

	// 创建用户服务实例，传入用户仓库
	userService := service.NewUserService(repository)

	// 创建用户处理器实例，传入用户服务
	userHandler := handler.NewUserHandler(userService)

	// 返回用户处理器实例
	return userHandler
}

package bootstrap

import (
	"go-my-blog/internal/handler"

	"gorm.io/gorm"
)

type Modules struct {
	UserHandler *handler.UserHandler
	PostHandler *handler.PostHandler
}

// initAllModules 初始化所有模块的函数
// 参数:
//   - db: 数据库连接对象，使用GORM库进行数据库操作
//
// 返回值:
//   - *Modules: 指向Modules结构体的指针，包含所有初始化后的模块
func InitAllModules(db *gorm.DB) *Modules {
	return &Modules{
		UserHandler: InitUserModule(db), // 初始化用户模块，并将返回的用户处理器赋值给Modules结构体的UserHandler字段
		PostHandler: InitPostModule(db), // 初始化文章模块，并将返回的文章处理器赋值给Modules结构体的PostHandler字段
	}
}

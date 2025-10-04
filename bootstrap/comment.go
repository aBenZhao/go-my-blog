package bootstrap

import (
	"go-my-blog/internal/handler"
	"go-my-blog/internal/repo"
	"go-my-blog/internal/service"

	"gorm.io/gorm"
)

// InitCommentModule 初始化评论模块，创建并返回一个CommentHandler实例
// 参数:
//   - db: 数据库连接对象，用于数据库操作
//
// 返回值:
//   - *handler.CommentHandler: 评论处理器实例，用于处理评论相关的请求
func InitCommentModule(db *gorm.DB) *handler.CommentHandler {

	// 创建评论数据仓库实例，用于数据库操作
	repository := repo.NewCommentRepository(db)

	// 创建评论服务实例，用于处理评论相关的业务逻辑
	commentService := service.NewCommentService(repository)

	// 创建评论处理器实例，用于处理HTTP请求和响应
	commentHandler := handler.NewCommentHandler(commentService)

	// 返回初始化完成的评论处理器
	return commentHandler
}

package bootstrap

import (
	"go-my-blog/internal/handler"
	"go-my-blog/internal/repo"
	"go-my-blog/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB

	// 仓库层
	UserRepo    *repo.UserRepository
	PostRepo    *repo.PostRepository
	CommentRepo *repo.CommentRepository

	// 服务层
	UserService    *service.UserSevice
	PostService    *service.PostService
	CommentService *service.CommentService

	// 处理器层
	UserHandler    *handler.UserHandler
	PostHandler    *handler.PostHandler
	CommentHandler *handler.CommentHandler
}

func NewContainer(db *gorm.DB) *Container {
	c := &Container{DB: db}

	// 初始化仓库层
	c.UserRepo = repo.NewUserRepository(db)
	c.PostRepo = repo.NewPostRepository(db)
	c.CommentRepo = repo.NewCommentRepository(db)

	// 初始化服务层
	c.UserService = service.NewUserService(c.UserRepo)
	c.PostService = service.NewPostService(c.PostRepo, c.UserRepo, c.CommentRepo)
	c.CommentService = service.NewCommentService(c.CommentRepo, c.UserRepo, c.PostRepo)

	// 初始化处理器层
	c.UserHandler = handler.NewUserHandler(c.UserService)
	c.PostHandler = handler.NewPostHandler(c.PostService)
	c.CommentHandler = handler.NewCommentHandler(c.CommentService)

	return c
}

func InitAllModules(db *gorm.DB) *Container {
	return NewContainer(db)
}

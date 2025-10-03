package bootstrap

import (
	"go-my-blog/internal/handler"
	"go-my-blog/internal/repo"
	"go-my-blog/internal/service"

	"gorm.io/gorm"
)

func InitPostModule(db *gorm.DB, userRepository *repo.UserRepository) *handler.PostHandler {
	postRepository := repo.NewPostRepository(db)

	postService := service.NewPostService(postRepository, userRepository)

	return handler.NewPostHandler(postService)
}

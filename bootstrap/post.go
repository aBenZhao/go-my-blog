package bootstrap

import (
	"go-my-blog/internal/handler"
	"go-my-blog/internal/repo"
	"go-my-blog/internal/service"

	"gorm.io/gorm"
)

func InitPostModule(db *gorm.DB) *handler.PostHandler {
	postRepository := repo.NewPostRepository(db)

	postService := service.NewPostService(postRepository, nil, nil)

	return handler.NewPostHandler(postService)
}

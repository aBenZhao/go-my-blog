package bootstrap

import (
	"go-my-blog/internal/handler"
	"go-my-blog/internal/repo"
	"go-my-blog/internal/service"

	"gorm.io/gorm"
)

func InitPostModule(db *gorm.DB) *handler.PostHandler {
	repository := repo.NewPostRepository(db)

	postService := service.NewPostService(repository)

	return handler.NewPostHandler(postService)
}

package repo

import (
	"go-my-blog/internal/model"
	"go-my-blog/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (pr *PostRepository) Create(post *model.Post) (*model.Post, error) {
	if err := pr.db.Create(post).Error; err != nil {
		logger.Error("PostRepository.Create db.Create is error", zap.Error(err))
		return nil, err
	}

	return post, nil
}

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
	if err := pr.db.Model(&model.Post{}).Create(post).Error; err != nil {
		logger.Error("PostRepository.Create db.Create is error", zap.Error(err))
		return nil, err
	}

	return post, nil
}

func (pr *PostRepository) GetById(id uint) (*model.Post, error) {
	var post model.Post
	if err := pr.db.Model(&model.Post{}).Where("id = ?", id).First(&post).Error; err != nil {
		logger.Error("PostRepository.GetById db.First is error", zap.Error(err))
		return nil, err
	}
	return &post, nil
}

func (pr *PostRepository) Updates(id uint, updateMap *map[string]interface{}) error {
	tx := pr.db.Model(&model.Post{}).Where("id = ?", id).Updates(updateMap)
	if tx.Error != nil {
		logger.Error("PostRepository.Updates db.Updates is error", zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		logger.Error("PostRepository.Updates db.Updates is error", zap.Error(tx.Error))
		return tx.Error
	}
	return nil

}

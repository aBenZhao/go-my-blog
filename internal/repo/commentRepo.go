package repo

import (
	"go-my-blog/internal/DTO"
	"go-my-blog/internal/model"
	"go-my-blog/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (cr CommentRepository) ListComments(dto DTO.ListCommentDTO) (*[]model.Comment, int64, error) {
	tx := cr.db.Model(&model.Comment{})

	if dto.Keyword != "" {
		tx = tx.Where("title LIKE ? OR content LIKE ?", "%"+dto.Keyword+"%", "%"+dto.Keyword+"%")
	}

	var total int64
	if err := tx.Where("post_id = ?", dto.PostId).Count(&total).Error; err != nil {
		logger.Error("PostRepository.ListPosts db.Count is error", zap.Error(err))
		return nil, 0, err
	}

	var comments []model.Comment
	offset := (dto.PageNum - 1) * dto.PageSize
	if err := tx.Where("post_id = ?", dto.PostId).Offset(offset).Limit(dto.PageSize).Find(&comments).Error; err != nil {
		logger.Error("PostRepository.ListPosts db.Find is error", zap.Error(err))
		return nil, 0, err
	}

	return &comments, total, nil
}

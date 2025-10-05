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

func (cr CommentRepository) DeleteByPostId(postId uint) error {
	tx := cr.db.Model(&model.Comment{}).Where("post_id = ?", postId).Delete(&model.Comment{})
	if tx.Error != nil {
		logger.Error("CommentRepository.DeleteByPostId is error", zap.Error(tx.Error))
		return tx.Error
	}
	return nil
}

func (cr CommentRepository) GetById(id uint) (*model.Comment, error) {
	var comment model.Comment
	tx := cr.db.Model(&model.Comment{}).Where("id = ?", id).First(&comment)
	if tx.Error != nil {
		logger.Error("CommentRepository.GetById is error", zap.Error(tx.Error))
		return nil, tx.Error
	}
	return &comment, nil
}

func (cr CommentRepository) DeleteById(id uint) error {
	tx := cr.db.Model(&model.Comment{}).Where("id = ?", id).Delete(&model.Comment{})
	if tx.Error != nil {
		logger.Error("CommentRepository.DeleteById is error", zap.Error(tx.Error))
		return tx.Error
	}
	return nil
}

func (cr CommentRepository) Create(comment *model.Comment) (*model.Comment, error) {
	if err := cr.db.Create(comment).Error; err != nil {
		logger.Error("CommentRepository.Create is error", zap.Error(err))
		return nil, err
	}
	return comment, nil
}

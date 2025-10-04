package repo

import (
	"go-my-blog/internal/DTO"
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

// GetById 根据ID获取文章信息的方法
// 参数:
//   - id: 文章的ID，类型为uint
//
// 返回值:
//   - *model.Post: 指向文章模型的指针，如果找到则返回文章数据
//   - error: 错误信息，如果查询过程中发生错误则返回错误
func (pr *PostRepository) GetById(id uint) (*model.Post, error) {
	// 声明一个Post结构体变量，用于存储查询结果
	var post model.Post
	// 执行数据库查询，根据ID查找文章
	// 如果查询过程中发生错误，记录错误日志并返回错误
	if err := pr.db.Model(&model.Post{}).Where("id = ?", id).First(&post).Error; err != nil {
		logger.Error("PostRepository.GetById db.First is error", zap.Error(err))
		return nil, err
	}
	// 成功查询后，返回文章数据的指针和nil错误
	return &post, nil
}

// Updates 更新指定ID的帖子信息
// 参数:
//
//	id - 帖子的ID
//	updateMap - 包含要更新字段的键值对映射
//
// 返回值:
//
//	error - 操作过程中遇到的错误，如果没有错误则返回nil
func (pr *PostRepository) Updates(id uint, updateMap *map[string]interface{}) error {
	// 创建数据库事务，更新指定ID的帖子记录
	tx := pr.db.Model(&model.Post{}).Where("id = ?", id).Updates(updateMap)
	// 检查数据库操作是否出错
	if tx.Error != nil {
		logger.Error("PostRepository.Updates db.Updates is error", zap.Error(tx.Error))
		return tx.Error
	}
	// 检查是否有行被影响，即是否有记录被更新
	if tx.RowsAffected == 0 {
		logger.Error("PostRepository.Updates db.Updates is error", zap.Error(tx.Error))
		return tx.Error
	}
	return nil

}

// Delete 删除指定ID的帖子
// 参数:
//
//	id: 要删除的帖子的ID
//
// 返回值:
//
//	error: 如果删除失败则返回错误信息
func (pr *PostRepository) Delete(id uint) error {
	// 使用GORM的Model方法和Where条件找到指定ID的帖子记录
	// 然后调用Delete方法删除该记录
	// 如果删除过程中出现错误，则返回该错误
	return pr.db.Model(&model.Post{}).Where("id = ?", id).Delete(&model.Post{}).Error
}

func (pr *PostRepository) ListPosts(dto *DTO.ListPostDTO) (*[]model.Post, int64, error) {
	tx := pr.db.Model(&model.Post{})

	if dto.Keyword != "" {
		tx = tx.Where("title LIKE ? OR content LIKE ?", "%"+dto.Keyword+"%", "%"+dto.Keyword+"%")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		logger.Error("PostRepository.ListPosts db.Count is error", zap.Error(err))
		return nil, 0, err
	}

	var posts []model.Post
	offset := (dto.PageNum - 1) * dto.PageSize
	if err := tx.Offset(offset).Limit(dto.PageSize).Find(&posts).Error; err != nil {
		logger.Error("PostRepository.ListPosts db.Find is error", zap.Error(err))
		return nil, 0, err
	}
	return &posts, total, nil
}

package repo

import (
	"errors"
	"go-my-blog/internal/model"
	"go-my-blog/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// UserRepository 的 UserRegister 方法用于处理用户注册
// 参数:
//   - user: 指向 model.User 结构体的指针，包含用户注册所需的信息
//
// 返回值:
//   - *model.User: 注册成功的用户信息
//   - error: 错误信息，如果注册失败则返回相应的错误
func (ur *UserRepository) UserRegister(user *model.User) (*model.User, error) {
	// 获取数据库连接
	db := ur.db
	// 尝试在数据库中创建新用户记录
	tx := db.Create(user)
	// 检查是否有错误发生
	if err := tx.Error; err != nil {
		// 记录错误日志
		logger.Error("UserRepository.UserRegister db.Create is error", zap.Error(err))
		// 返回空用户指针和错误信息
		return nil, err
	}
	// 返回成功注册的用户信息和 nil 错误
	return user, nil
}

func (ur *UserRepository) FindByUserName(username string) (*model.User, error) {
	var user model.User
	// 获取数据库连接
	db := ur.db
	tx := db.Where("username = ?", username).First(&user)
	if tx.Error != nil {
		// 用户不存在
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			logger.Warn("UserRepository.FindByUserName user not found:"+username, zap.Error(tx.Error))
		}
		logger.Error("UserRepository.FindByUserName db.Where is error", zap.Error(tx.Error))
		return nil, tx.Error
	}
	return &user, nil
}

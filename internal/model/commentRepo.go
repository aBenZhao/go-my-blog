package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `gorm:"type:bigint;primaryKey;autoIncrement;comment:评论唯一标识" json:"id"`
	Content   string         `gorm:"type:text;not null;comment:评论内容" json:"content"`
	UserID    uint           `gorm:"type:bigint;not null;index:idx_comment_user;comment:评论者ID" json:"user_id"`
	PostID    uint           `gorm:"type:bigint;not null;index:idx_comment_post;comment:所属文章ID" json:"post_id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"comment:软删除标记" json:"-"`
	// 关联关系：评论的作者和所属文章
	User User `gorm:"foreignKey:UserID" json:"user"` // 预加载评论者信息
	Post Post `gorm:"foreignKey:PostID" json:"post"` // 可选：关联文章信息
}

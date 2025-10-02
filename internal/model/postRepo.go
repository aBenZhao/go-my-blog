package model

import (
	"time"

	"gorm.io/gorm"
)

// Post 文章模型
type Post struct {
	ID        uint           `gorm:"type:bigint;primaryKey;autoIncrement;comment:文章唯一标识" json:"id"`
	Title     string         `gorm:"type:varchar(200);not null;index:idx_title;comment:文章标题" json:"title"`
	Content   string         `gorm:"type:text;not null;comment:文章内容" json:"content"`
	UserID    uint           `gorm:"type:bigint;not null;index:idx_post_user;comment:作者ID" json:"user_id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"comment:软删除标记" json:"deleted_at"`
	// 关联作者：无需级联（删除文章不影响用户），保持不变
	User User `gorm:"foreignKey:UserID" json:"user"`
	// 关键修改：添加 OnDelete:CASCADE，与 SQL 级联删除逻辑对齐
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments"`
}

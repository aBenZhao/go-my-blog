package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"type:bigint;primaryKey;autoIncrement;comment:用户唯一标识" json:"id"`
	Username  string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_username;comment:用户名（唯一）" json:"username"`
	Password  string         `gorm:"type:varchar(100);not null;comment:加密存储的密码" json:"password"`
	Email     string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_email;comment:邮箱（唯一）" json:"email"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"comment:软删除标记" json:"deleted_at"`
	// 添加 OnDelete:CASCADE，与 SQL 级联删除逻辑对齐
	Posts    []Post    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"posts"`
	Comments []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"comments"`
}

//func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
//
//}

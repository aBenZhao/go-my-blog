package request

import "github.com/go-playground/validator/v10"

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" validate:"required"` // 用户名（必填）
	Password string `json:"password" validate:"required"` // 密码（必填）
}

// 校验参数
func (req *LoginRequest) Validate() error {
	return validator.New().Struct(req)
}

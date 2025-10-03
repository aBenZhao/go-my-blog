package handler

import (
	"go-my-blog/internal/DTO"
	"go-my-blog/internal/request"
	"go-my-blog/internal/response"
	"go-my-blog/internal/service"
	"go-my-blog/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserSevice
}

func NewUserHandler(userService *service.UserSevice) *UserHandler {
	return &UserHandler{userService}
}

// UserRegister 处理用户注册请求
// UserRegister 处理用户注册请求的函数
// @Summary 用户注册
// @Description 接收用户注册请求，验证参数并处理注册逻辑
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "注册请求参数"
// @Success 200 {object} gin.H "注册成功，返回用户信息"
// @Failure 400 {object} gin.H "参数错误"
// @Failure 500 {object} gin.H "服务器内部错误"
// @Router /register [post]
func (uh *UserHandler) UserRegister(c *gin.Context) {
	// 1. 绑定请求参数（JSON 格式）
	// 从请求体中解析JSON数据到RegisterRequest结构体
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("注册参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误：" + err.Error()})
		return
	}

	// 2. 调用服务层处理注册逻辑
	// 将请求数据转换为DTO格式，调用用户服务层处理注册
	userRegisterDTO := DTO.UserRegisterDTO{Username: req.Username, Password: req.Password, Email: req.Email}
	user, err := uh.userService.UserRegister(&userRegisterDTO)
	if err != nil {
		logger.Error("注册业务处理失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "注册失败：" + err.Error()})
		return
	}

	// 使用copier将用户实体复制到响应DTO
	var userResponse response.UserResponse
	err = copier.Copy(&userResponse, &user)
	if err != nil {
		logger.Error("注册业务处理失败：响应对象复制失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "注册失败：" + err.Error()})
	}

	//3. 返回成功响应
	// 返回HTTP 200状态码和成功注册的用户信息
	c.JSON(http.StatusOK, gin.H{
		"msg":  "注册成功",
		"data": userResponse,
	})

}

func (uh *UserHandler) UserLogin(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("登入参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误：" + err.Error()})
		return
	}
	dto := DTO.LoginDTO{req.Username, req.Password}
	loginResponse, err := uh.userService.UserLogin(&dto)
	if err != nil {
		logger.Error("登录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "登录失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "登录成功",
		"data": loginResponse,
	})
}

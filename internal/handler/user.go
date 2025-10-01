package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserRegister 处理用户注册请求
func UserRegister(c *gin.Context) {
	//// 1. 绑定请求参数（JSON 格式）
	//var req model.RegisterRequest
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	logger.Warn("注册参数绑定失败", zap.Error(err))
	//	c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误：" + err.Error()})
	//	return
	//}
	//
	//// 2. 调用服务层处理注册逻辑
	//user, err := service.UserRegister(req.Username, req.Password, req.Email)
	//if err != nil {
	//	logger.Error("注册业务处理失败", zap.Error(err))
	//	c.JSON(http.StatusInternalServerError, gin.H{"msg": "注册失败：" + err.Error()})
	//	return
	//}
	//
	////3. 返回成功响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg":  "注册成功",
	//	"data": user,
	//})

	c.JSON(http.StatusOK, gin.H{
		"msg":  "注册成功",
		"data": "",
	})
}

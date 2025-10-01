package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-my-blog/pkg/jwt" // 引入上面实现的 JWT 工具类
)

// JWTAuth JWT 认证中间件：验证请求中的 Token 有效性，通过后将用户 ID 存入上下文
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头中获取 Token（格式：Authorization: Bearer <token>）
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请先登录（未携带 Authorization 头）",
			})
			c.Abort() // 拦截请求，不再执行后续 handler
			return
		}

		// 2. 校验 Token 格式（必须是 "Bearer <token>" 格式）
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token 格式错误（正确格式：Bearer <token>）",
			})
			c.Abort()
			return
		}

		// 3. 验证 Token 有效性并提取用户 ID
		userID, err := jwt.VerifyToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token 无效或已过期：" + err.Error(),
			})
			c.Abort()
			return
		}

		// 4. 将用户 ID 存入上下文，供后续 handler 使用（如创建文章时记录作者 ID）
		c.Set("userID", userID)

		// 5. 继续执行后续中间件或 handler
		c.Next()
	}
}

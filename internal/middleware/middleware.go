package middleware

import (
	"go-my-blog/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 自定义日志中间件：记录请求的方法、路径、状态码、耗时等信息
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 预处理：记录请求开始时间、请求方法、路径
		startTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		remoteAddr := c.ClientIP() // 客户端 IP

		// 2. 执行后续的中间件和 handler（必须调用，否则请求会被拦截）
		c.Next()

		// 3. 后处理：获取响应状态码、计算耗时，打印日志
		statusCode := c.Writer.Status()
		duration := time.Since(startTime)

		// 用自定义日志包打印请求日志（结构化字段便于分析）
		logger.Info(
			"请求访问日志",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("remote_addr", remoteAddr),
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration),
		)
	}
}

// GinRecovery 异常恢复中间件：捕获 handler 中的 panic，返回 500 错误，避免服务崩溃
// 参数：debug 是否在开发环境显示错误堆栈
func GinRecovery(debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil { // 捕获 panic
				// 1. 打印错误日志（包含堆栈信息）
				logger.Error(
					"请求处理 panic 恢复",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.Stack("stack_trace"), // 记录堆栈信息，便于排查问题
				)

				// 2. 返回 500 响应（开发环境显示错误详情，生产环境隐藏）
				if debug {
					// 开发环境：返回错误信息和堆栈（方便调试）
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg":   "服务器内部错误",
						"error": err,
						"stack": zap.Stack("").String, // 堆栈信息
					})
				} else {
					// 生产环境：隐藏敏感错误信息
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "服务器内部错误，请稍后重试",
					})
				}

				// 3. 终止请求链，不再执行后续逻辑
				c.Abort()
			}
		}()

		// 执行后续中间件和 handler
		c.Next()
	}
}

// Cors 跨域处理中间件：解决前端调用 API 时的跨域限制（如浏览器的同源策略）
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 设置允许跨域的源（* 表示允许所有源，生产环境可指定具体域名如 https://your-frontend.com）
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*") // 无 Origin 时允许所有
		}

		// 2. 允许的请求方法（GET/POST/PUT/DELETE 等）
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 3. 允许的请求头（包含自定义头如 Authorization）
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 4. 允许前端读取的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin")

		// 5. 是否允许携带 Cookie（跨域请求时）
		c.Header("Access-Control-Allow-Credentials", "true")

		// 6. 处理 OPTIONS 预检请求（浏览器会先发送 OPTIONS 检查跨域权限）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 204 表示允许跨域
			return
		}

		// 执行后续中间件和 handler
		c.Next()
	}
}

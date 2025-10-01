package router

import (
	"go-my-blog/config/priority_config"
	"go-my-blog/internal/handler"    // 引入处理器（接口实现逻辑）
	"go-my-blog/internal/middleware" // 引入中间件（如认证、日志）

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由：将所有 API 注册到 Gin 引擎
func InitRouter(r *gin.Engine) {
	// 1. 全局中间件：所有路由都会经过的中间件（如日志、跨域）
	r.Use(middleware.GinLogger())                                         // 自定义日志中间件（记录请求日志）
	r.Use(middleware.GinRecovery(priority_config.PriorityConf.Gin.Debug)) // 异常恢复中间件（避免服务因 panic 崩溃）
	r.Use(middleware.Cors())                                              // 跨域处理中间件（前端调用 API 时需要）

	// 2. 无需认证的路由组（公开接口）
	public := r.Group("/api/v1")
	{
		// 用户相关公开接口
		public.POST("/register", handler.UserRegister) // 用户注册
		//public.POST("/login", handler.UserLogin)       // 用户登录

		// 文章相关公开接口
		//public.GET("/posts", handler.PostList)       // 文章列表（分页）
		//public.GET("/posts/:id", handler.PostDetail) // 文章详情

		// 评论相关公开接口
		//public.GET("/posts/:postID/comments", handler.CommentList) // 文章的评论列表
	}

	// 3. 需要认证的路由组（需登录才能访问）
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth()) // JWT 认证中间件：验证 token 有效性
	{
		// 用户相关私有接口（需登录）
		//auth.GET("/user/profile", handler.UserProfile)   // 获取当前用户信息
		//auth.PUT("/user/profile", handler.UpdateProfile) // 更新用户信息

		// 文章相关私有接口（需登录）
		//auth.POST("/posts", handler.CreatePost)       // 创建文章
		//auth.PUT("/posts/:id", handler.UpdatePost)    // 更新文章
		//auth.DELETE("/posts/:id", handler.DeletePost) // 删除文章

		// 评论相关私有接口（需登录）
		//auth.POST("/posts/:postID/comments", handler.CreateComment) // 发布评论
		//auth.DELETE("/comments/:id", handler.DeleteComment)         // 删除自己的评论
	}
}

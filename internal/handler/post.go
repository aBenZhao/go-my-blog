package handler

import (
	"go-my-blog/internal/DTO"
	"go-my-blog/internal/request"
	"go-my-blog/internal/response"
	"go-my-blog/internal/service"
	"go-my-blog/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type PostHandler struct {
	PostService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService}
}

// CreatePost 是处理创建文章请求的方法
// 它从上下文中获取用户ID，验证请求参数，并调用服务层创建文章
func (ph *PostHandler) CreatePost(c *gin.Context) {
	// 从上下文中获取用户ID，检查是否存在
	userID, exists := c.Get("userID")
	if !exists {
		logger.Warn("用户授权失败")
		// 如果用户ID不存在，返回未授权错误
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized: no user ID found in context"})
		return
	}

	// 声明创建文章的请求结构体，并绑定请求体数据
	var req request.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果参数绑定失败，返回参数错误信息
		logger.Error("创建参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误：" + err.Error()})
		return
	}

	// 与CreatePostRequest中的validate标签配合使用，需要在此显示校验
	if err := validator.New().Struct(req); err != nil {
		logger.Error("创建参数校验失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误：" + err.Error()})
	}

	// 调用PostService的CreatePost方法创建文章，传入用户ID和请求参数
	createPostDTO := DTO.CreatePostDTO{Title: req.Title, Content: req.Content}
	postRespDTO, err := ph.PostService.CreatePost(userID.(uint), &createPostDTO)
	if err != nil {
		logger.Error("创建文章失败", zap.Error(err))
		// 如果创建失败，返回服务器内部错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建文章失败：" + err.Error()})
		return
	}

	var postResp response.CreatePostResponse
	if err = copier.Copy(&postResp, &postRespDTO); err != nil {
		logger.Error("拷贝失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "注册失败：" + err.Error()})
	}

	// 创建成功，返回成功响应和文章数据
	c.JSON(http.StatusOK, gin.H{"msg": "文章创建成功", "data": postResp})
}

func (ph *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		logger.Error("文章ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID不能为空"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.Warn("用户授权失败")
		// 如果用户ID不存在，返回未授权错误
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized: no user ID found in context"})
		return
	}

	var req request.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新文章参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数错误：" + err.Error()})
		return
	}

	var updatePostDTO DTO.UpdatePostDTO
	if err := copier.Copy(&updatePostDTO, &req); err != nil {
		logger.Error("拷贝失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "拷贝失败：" + err.Error()})
		return
	}
	idUint, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		logger.Error("文章ID格式错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID格式错误：" + err.Error()})
		return
	}
	updatePostDTO.ID = uint(idUint)
	updatePostDTO.UserID = userID.(uint)

	postDTO, err := ph.PostService.UpdatePost(&updatePostDTO)
	if err != nil {
		logger.Error("更新文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新文章失败：" + err.Error()})
		return
	}

	var updatePostResponse response.UpdatePostResponse
	if err := copier.Copy(&updatePostResponse, &postDTO); err != nil {
		logger.Error("拷贝失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "拷贝失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "更新文章成功", "data": updatePostResponse})

}

func (ph *PostHandler) DeletePost(context *gin.Context) {
	idStr := context.Param("id")
	if idStr == "" {
		logger.Error("文章ID不能为空")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID不能为空"})
		return
	}

	userID, exists := context.Get("userID")
	if !exists {
		logger.Warn("用户授权失败")
		// 如果用户ID不存在，返回未授权错误
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized: no user ID found in context"})
		return
	}

	idUint, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		logger.Error("文章ID格式错误", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID格式错误：" + err.Error()})
		return
	}

	err = ph.PostService.DeletePost(uint(idUint), userID.(uint))
	if err != nil {
		logger.Error("删除文章失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "删除文章失败：" + err.Error()})
	}
	context.JSON(http.StatusOK, gin.H{"msg": "删除文章成功"})
}

func (ph *PostHandler) PostList(context *gin.Context) {
	var req request.PostListRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		logger.Error("获取文章列表参数绑定失败", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数错误：" + err.Error()})
		return
	}
	req.SetDefault()

	userID, exists := context.Get("userID")
	if !exists {
		logger.Warn("用户授权失败")
		// 如果用户ID不存在，返回未授权错误
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized: no user ID found in context"})
		return
	}

	var listPostDTO DTO.ListPostDTO
	copyErr := copier.Copy(&listPostDTO, &req)
	if copyErr != nil {
		logger.Error("拷贝失败", zap.Error(copyErr))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "拷贝失败：" + copyErr.Error()})
		return
	}
	listPostDTO.UserID = userID.(uint)

	postDTOList, err := ph.PostService.PostList(&listPostDTO)
	if err != nil {
		logger.Error("获取文章列表失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "获取文章列表失败：" + err.Error()})
		return
	}

	var postListResponse response.PostListResponse
	if err := copier.Copy(&postListResponse, &postDTOList); err != nil {
		logger.Error("拷贝失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "拷贝失败：" + err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"msg": "获取文章列表成功", "data": postListResponse})
}

func (ph *PostHandler) PostDetail(context *gin.Context) {
	idStr := context.Param("id")
	if idStr == "" {
		logger.Error("文章ID不能为空")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID不能为空"})
		return
	}

	parseUint, parseErr := strconv.ParseUint(idStr, 10, 0)
	if parseErr != nil {
		logger.Error("文章ID格式错误", zap.Error(parseErr))
		context.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID格式错误：" + parseErr.Error()})
	}

	postDetailDTO, err := ph.PostService.PostDetail(parseUint)
	if err != nil {
		logger.Error("获取文章详情失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "获取文章详情失败：" + err.Error()})
		return
	}

	var postResp response.PostDetailResponse
	if err := copier.Copy(&postResp, &postDetailDTO); err != nil {
		logger.Error("拷贝失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "拷贝失败：" + err.Error()})
	}

	context.JSON(http.StatusOK, gin.H{"msg": "获取文章详情成功", "data": postResp})
}

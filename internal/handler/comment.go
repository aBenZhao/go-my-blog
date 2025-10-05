package handler

import (
	"go-my-blog/internal/DTO"
	"go-my-blog/internal/repo"
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

type CommentHandler struct {
	commentService *service.CommentService
}

func (h CommentHandler) GetCommentRepo() *repo.CommentRepository {
	return h.commentService.GetCommentRepo()
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (ch CommentHandler) DeleteComment(context *gin.Context) {
	commentIdStr := context.Param("id")
	if commentIdStr == "" {
		logger.Error("评论id不能为空")
		context.JSON(400, gin.H{"error": "commentId is required"})
	}

	userID, exists := context.Get("userID")
	if !exists {
		logger.Warn("用户授权失败")
		// 如果用户ID不存在，返回未授权错误
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized: no user ID found in context"})
		return
	}

	commentId, err := strconv.ParseUint(commentIdStr, 10, 0)
	if err != nil {
		logger.Error("评论ID格式错误", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"msg": "评论ID格式错误：" + err.Error()})
		return
	}

	err = ch.commentService.DeleteComment(uint(commentId), userID.(uint))
	if err != nil {
		logger.Error("删除评论失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "删除评论失败：" + err.Error()})
	}
	context.JSON(http.StatusOK, gin.H{"msg": "删除评论成功"})
}

func (ch CommentHandler) CreateComment(context *gin.Context) {

	userID, exists := context.Get("userID")
	if !exists {
		logger.Warn("用户授权失败")
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized: no user ID found in context"})
		return
	}
	postIDStr := context.Param("postID")
	if postIDStr == "" {
		logger.Error("文章ID不能为空")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID不能为空"})
		return
	}
	postID, err := strconv.ParseUint(postIDStr, 10, 0)
	if err != nil {
		logger.Error("文章ID格式错误", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID格式错误：" + err.Error()})
		return
	}

	var req request.CreateCommentRequest
	req.PostID = uint(postID)
	if err := context.ShouldBindJSON(&req); err != nil {
		logger.Error("创建参数绑定失败", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误：" + err.Error()})
		return
	}

	if err := validator.New().Struct(req); err != nil {
		logger.Error("创建参数校验失败", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误：" + err.Error()})
	}

	createCommentDTO := DTO.CreateCommentDTO{Content: req.Content, PostID: req.PostID}
	commentRespDTO, err := ch.commentService.CreateComment(userID.(uint), &createCommentDTO)
	if err != nil {
		logger.Error("创建评论失败", zap.Error(err))
		// 如果创建失败，返回服务器内部错误信息
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "创建评论失败：" + err.Error()})
		return
	}

	var commentResp response.CreateCommentResponse
	if err = copier.Copy(&commentResp, &commentRespDTO); err != nil {
		logger.Error("拷贝失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "拷贝失败：" + err.Error()})
	}

	context.JSON(http.StatusOK, gin.H{"msg": "评论创建成功", "data": commentResp})
}

func (ch CommentHandler) CommentList(context *gin.Context) {
	postIDStr := context.Param("postID")
	if postIDStr == "" {
		logger.Error("文章ID不能为空")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID不能为空"})
		return
	}
	postID, err := strconv.ParseUint(postIDStr, 10, 0)
	if err != nil {
		logger.Error("文章ID格式错误", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"msg": "文章ID格式错误：" + err.Error()})
		return
	}

	comments, err := ch.commentService.CommentList(uint(postID))
	if err != nil {
		logger.Error("获取评论列表失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "获取评论列表失败：" + err.Error()})
	}
	context.JSON(http.StatusOK, gin.H{"msg": "获取评论列表成功", "data": comments})
}

package handler

import (
	"go-my-blog/internal/repo"
	"go-my-blog/internal/service"
	"go-my-blog/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

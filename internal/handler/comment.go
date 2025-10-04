package handler

import (
	"go-my-blog/internal/repo"
	"go-my-blog/internal/service"
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

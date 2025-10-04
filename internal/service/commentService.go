package service

import "go-my-blog/internal/repo"

type CommentService struct {
	commentRepo *repo.CommentRepository
}

func (s CommentService) GetCommentRepo() *repo.CommentRepository {
	return s.commentRepo
}

func NewCommentService(commentRepo *repo.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

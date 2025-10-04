package service

import (
	"database/sql"
	"errors"
	"go-my-blog/internal/repo"
	"go-my-blog/pkg/logger"

	"go.uber.org/zap"
)

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

func (cs CommentService) DeleteComment(commentId uint, userId uint) error {
	// 根据ID从数据库中获取文章信息
	comment, err := cs.commentRepo.GetById(commentId)
	if err != nil {
		// 检查错误是否为"未找到行"的错误
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("评论不存在", zap.Error(err))
			return err
		}
		// 其他数据库错误处理
		logger.Error("评论查询失败", zap.Error(err))
		return err
	}

	// 验证请求删除的用户是否为评论作者
	if userId != comment.UserID {
		logger.Error("登录用户非评论作者，不允许删除评论")
		return errors.New("登录用户非评论作者，不允许删除评论")
	}

	err = cs.commentRepo.DeleteById(commentId)
	if err != nil {
		logger.Error("删除评论异常", zap.Error(err))
		return err
	}
	return nil
}

package service

import (
	"database/sql"
	"errors"
	"go-my-blog/internal/DTO"
	"go-my-blog/internal/model"
	"go-my-blog/internal/repo"
	"go-my-blog/pkg/logger"

	"go.uber.org/zap"
)

type CommentService struct {
	commentRepo *repo.CommentRepository
	userRepo    *repo.UserRepository
	postRepo    *repo.PostRepository
}

func (s CommentService) GetCommentRepo() *repo.CommentRepository {
	return s.commentRepo
}

func NewCommentService(commentRepo *repo.CommentRepository, userRepo *repo.UserRepository, postRepo *repo.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		postRepo:    postRepo,
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

func (cs CommentService) CreateComment(userID uint, d *DTO.CreateCommentDTO) (*DTO.CreateCommentDTO, error) {
	_, err := cs.postRepo.GetById(d.PostID)
	if err != nil {
		logger.Error("文章不存在", zap.Error(err))
		return nil, err
	}
	var comment = &model.Comment{PostID: d.PostID, UserID: userID, Content: d.Content}
	commentResult, err := cs.commentRepo.Create(comment)
	if err != nil {
		logger.Error("评论创建失败", zap.Error(err))
		return nil, err
	}

	var resultDTO = &DTO.CreateCommentDTO{PostID: commentResult.PostID, Content: commentResult.Content}

	return resultDTO, nil

}

func (cs CommentService) CommentList(postId uint) (*[]DTO.CommentDetailDTO, error) {
	_, err := cs.postRepo.GetById(postId)
	if err != nil {
		logger.Error("文章不存在", zap.Error(err))
		return nil, err
	}

	comments, _, err := cs.commentRepo.ListComments(DTO.ListCommentDTO{PostId: postId})
	if err != nil {
		logger.Error("评论列表查询失败", zap.Error(err))
		return nil, err
	}

	var commentDetailsDTO []DTO.CommentDetailDTO
	for _, comment := range *comments {
		var commentDetailDTO DTO.CommentDetailDTO
		commentDetailDTO.ID = comment.ID
		commentDetailDTO.Content = comment.Content
		commentDetailDTO.UserID = comment.UserID
		commentDetailDTO.CreatedAt = comment.CreatedAt.Format("2006-01-02 15:04:05")
		commentDetailDTO.UpdatedAt = comment.UpdatedAt.Format("2006-01-02 15:04:05")
		commentDetailsDTO = append(commentDetailsDTO, commentDetailDTO)
	}

	return &commentDetailsDTO, nil
}

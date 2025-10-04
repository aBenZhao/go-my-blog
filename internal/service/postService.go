package service

import (
	"database/sql"
	"errors"
	"go-my-blog/internal/DTO"
	"go-my-blog/internal/model"
	"go-my-blog/internal/repo"
	"go-my-blog/pkg/logger"
	"time"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type PostService struct {
	PostRepo    *repo.PostRepository
	UserRepo    *repo.UserRepository
	CommentRepo *repo.CommentRepository
}

func NewPostService(postRepo *repo.PostRepository, userRepo *repo.UserRepository, commentRepo *repo.CommentRepository) *PostService {
	return &PostService{
		PostRepo:    postRepo,
		UserRepo:    userRepo,
		CommentRepo: commentRepo,
	}
}

func (ps *PostService) CreatePost(userID uint, createPostDTO *DTO.CreatePostDTO) (*DTO.CreatePostDTO, error) {
	var post model.Post
	if err := copier.Copy(&post, createPostDTO); err != nil {
		logger.Error("PostService.CreatePost copier.Copy is error!", zap.Error(err))
		return nil, err
	}
	post.UserID = userID

	postResp, err := ps.PostRepo.Create(&post)
	if err != nil {
		logger.Error("PostService.CreatePost PostRepo.Create is error!", zap.Error(err))
		return nil, err
	}

	var postResult DTO.CreatePostDTO
	if err := copier.Copy(&postResult, &postResp); err != nil {
		logger.Error("PostService.CreatePost copier.Copy is error!", zap.Error(err))
		return nil, err
	}
	return &postResult, nil
}

func (ps *PostService) UpdatePost(updatePostDTO *DTO.UpdatePostDTO) (*DTO.UpdatePostDTO, error) {
	id := updatePostDTO.ID
	post, err := ps.PostRepo.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("文章不存在", zap.Error(err))
			return nil, err
		}
		logger.Error("文章查询失败", zap.Error(err))
		return nil, err
	}
	user, err := ps.UserRepo.FindById(updatePostDTO.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("用户不存在", zap.Error(err))
			return nil, err
		}
		logger.Error("用户查询失败", zap.Error(err))
	}

	if post.UserID != user.ID {
		logger.Error("登录用户非文章作者，不允许更新文章")
		return nil, errors.New("登录用户非文章作者，不允许更新文章")
	}

	updateMap := make(map[string]interface{})
	updateMap["title"] = updatePostDTO.Title
	updateMap["content"] = updatePostDTO.Content
	updateMap["updated_at"] = time.Now()
	if ps.PostRepo.Updates(id, &updateMap) != nil {
		logger.Error("文章更新失败", zap.Error(err))
		return nil, err
	}

	var updateAffectedPostDTO DTO.UpdatePostDTO
	updateAffectedPost, _ := ps.PostRepo.GetById(id)
	if copier.Copy(&updateAffectedPostDTO, &updateAffectedPost) != nil {
		logger.Error("PostService.UpdatePost copier.Copy is error!", zap.Error(err))
		return nil, err
	}

	return &updateAffectedPostDTO, nil
}

// DeletePost 删除文章的方法
// 参数:
//
//	id: 要删除的文章ID
//	userId: 请求删除文章的用户ID
//
// 返回值:
//
//	error: 操作过程中遇到的错误，如果删除成功则返回nil
func (ps *PostService) DeletePost(id uint, userId uint) error {
	// 根据ID从数据库中获取文章信息
	post, err := ps.PostRepo.GetById(id)
	if err != nil {
		// 检查错误是否为"未找到行"的错误
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("文章不存在", zap.Error(err))
			return err
		}
		// 其他数据库错误处理
		logger.Error("文章查询失败", zap.Error(err))
		return err
	}

	// 验证请求删除的用户是否为文章作者
	if userId != post.UserID {
		logger.Error("登录用户非文章作者，不允许删除文章")
		return errors.New("登录用户非文章作者，不允许删除文章")
	}

	// 调用仓储层执行删除操作
	err = ps.PostRepo.Delete(id)

	if err != nil {
		logger.Error("文章删除失败", zap.Error(err))
		return err
	}

	err = ps.CommentRepo.DeleteByPostId(post.ID)
	if err != nil {
		logger.Error("文章删除失败", zap.Error(err))
		return err
	}
	return nil
}

func (ps *PostService) PostList(listPostDTO *DTO.ListPostDTO) (*DTO.PostListDTO, error) {
	posts, total, err := ps.PostRepo.ListPosts(listPostDTO)
	if err != nil {
		logger.Error("PostService.PostList PostRepo.ListPosts is error!", zap.Error(err))
		return nil, err
	}

	var postDTO []DTO.PostDTO
	if err := copier.Copy(&postDTO, &posts); err != nil {
		logger.Error("PostService.PostList copier.Copy is error!", zap.Error(err))
		return nil, err
	}
	var postListDTO DTO.PostListDTO
	postListDTO.Posts = postDTO
	postListDTO.Total = total
	postListDTO.PageNum = listPostDTO.PageNum
	postListDTO.PageSize = listPostDTO.PageSize
	return &postListDTO, nil

}

func (ps *PostService) PostDetail(postId uint64) (*DTO.PostDetailDTO, error) {
	post, err := ps.PostRepo.GetById(uint(postId))
	if err != nil {
		logger.Error("PostService.PostDetail PostRepo.GetById is error!", zap.Error(err))
		return nil, err
	}

	user, err := ps.UserRepo.FindById(post.UserID)
	if err != nil {
		logger.Error("PostService.PostDetail UserRepo.FindById is error!", zap.Error(err))
		return nil, err
	}

	var listCommentDTO DTO.ListCommentDTO
	listCommentDTO.PostId = postId
	listCommentDTO.PageNum = 1
	listCommentDTO.PageSize = 10
	comments, _, err := ps.CommentRepo.ListComments(listCommentDTO)
	if err != nil {
		logger.Error("PostService.PostDetail CommentRepo.ListComments is error!", zap.Error(err))
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

	var postDetailDTO DTO.PostDetailDTO
	postDetailDTO.ID = post.ID
	postDetailDTO.Title = post.Title
	postDetailDTO.Content = post.Content
	postDetailDTO.CreatedAt = post.CreatedAt.Format("2006-01-02 15:04:05")
	postDetailDTO.UpdatedAt = post.UpdatedAt.Format("2006-01-02 15:04:05")
	postDetailDTO.UserID = post.UserID
	postDetailDTO.Username = user.Username
	postDetailDTO.Comments = commentDetailsDTO

	return &postDetailDTO, nil

}

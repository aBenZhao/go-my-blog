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
	PostRepo *repo.PostRepository
	UserRepo *repo.UserRepository
}

func NewPostService(postRepo *repo.PostRepository, userRepo *repo.UserRepository) *PostService {
	return &PostService{
		PostRepo: postRepo,
		UserRepo: userRepo,
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

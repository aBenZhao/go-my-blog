package service

import (
	"go-my-blog/internal/DTO"
	"go-my-blog/internal/model"
	"go-my-blog/internal/repo"
	"go-my-blog/pkg/logger"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type PostService struct {
	PostRepo *repo.PostRepository
}

func NewPostService(postRepo *repo.PostRepository) *PostService {
	return &PostService{
		PostRepo: postRepo,
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

package DTO

type CommentDetailDTO struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
	//Username  string    `json:"username"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListCommentDTO struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	UserID   uint   `json:"user_id"`
	PostId   uint   `json:"post_id"`
}

type CreateCommentDTO struct {
	Content string `json:"content"`
	PostID  uint   `json:"post_id"`
}

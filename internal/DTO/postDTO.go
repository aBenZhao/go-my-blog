package DTO

type CreatePostDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostDTO struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
}

type ListPostDTO struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	UserID   uint   `json:"user_id"`
}

type PostDTO struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostListDTO struct {
	Posts    []PostDTO `json:"posts"`
	Total    int64     `json:"total"`
	PageNum  int       `form:"page_num"`
	PageSize int       `form:"page_size"`
}

type PostDetailDTO struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserID    uint   `json:"userId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Username  string `json:"username"`

	Comments []CommentDetailDTO `json:"comments"`
}

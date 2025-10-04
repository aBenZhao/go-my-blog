package response

type CreatePostResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostListResponse struct {
	Posts    []PostResponse `json:"posts"`
	Total    int            `json:"total"`
	PageNum  int            `form:"page_num"`
	PageSize int            `form:"page_size"`
}

type PostDetailResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserID    uint   `json:"userId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Username  string `json:"username"`

	Comments []CommentDetailResponse `json:"comments"`
}

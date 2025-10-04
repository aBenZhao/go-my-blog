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
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Username string `json:"username"`
}

type PostListResponse struct {
	Posts    []PostResponse `json:"posts"`
	Total    int            `json:"total"`
	PageNum  int            `form:"page_num"`
	PageSize int            `form:"page_size"`
}

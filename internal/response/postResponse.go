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

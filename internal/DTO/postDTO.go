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

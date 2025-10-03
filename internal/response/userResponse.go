package response

type UserResponse struct {
	Username string `json:"username"`
	Password string `json:"password" copier:"-"`
	Email    string `json:"email"`
}

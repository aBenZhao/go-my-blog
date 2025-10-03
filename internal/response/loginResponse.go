package response

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
	ExpiresAt   int64  `json:"expires_at"` // 令牌过期时间（时间戳，单位秒）
}

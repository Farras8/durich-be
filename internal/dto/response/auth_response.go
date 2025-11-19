package response

type RegisterAdminResponse struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	RoleSystem string `json:"role_system"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
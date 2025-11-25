package response

import "durich-be/internal/domain"

type RegisterResponse struct {
	Email string `json:"email"`
}

type LoginResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	Roles        []domain.UserRole `json:"roles"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

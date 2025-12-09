package requests

import (
	"durich-be/internal/domain"

	"github.com/segmentio/ksuid"
)

// Authentication Request DTOs
type RegisterAdmin struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	LocationID string `json:"location_id"` // Optional: If empty, treated as Central Admin
}

type Login struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Internal Authentication DTOs
type UserAuth struct {
	AuthID         string            `json:"auth_id"`
	UserID         string            `json:"user_id"`
	Email          string            `json:"email"`
	Role           []domain.UserRole `json:"role"`
	LocationID     string            `json:"location_id"`
	RefreshTokenID string            `json:"refresh_token_id,omitempty"`
}

type CreateAuth struct {
	UserEmail         string `json:"user_email"`
	Password          string `json:"password"`
	EncryptedPassword string `json:"-"`
}

func (receiver CreateAuth) ToDomain() domain.Authentication {
	return domain.Authentication{
		UserEmail:         receiver.UserEmail,
		Password:          receiver.Password,
		EncryptedPassword: receiver.EncryptedPassword,
	}
}

func ToTokenPayload(record domain.Authentication) UserAuth {
	roles := []domain.UserRole{}
	userID := ""
	locationID := ""
	if record.User != nil {
		roles = record.User.RoleSystem
		userID = record.User.ID
		if record.User.CurrentLocationID != nil {
			locationID = *record.User.CurrentLocationID
		}
	}

	return UserAuth{
		AuthID:         record.ID,
		UserID:         userID,
		Email:          record.UserEmail,
		Role:           roles,
		LocationID:     locationID,
		RefreshTokenID: ksuid.New().String(),
	}
}

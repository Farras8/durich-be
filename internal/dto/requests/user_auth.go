package requests

import (
	"durich-be/internal/domain"

	"github.com/segmentio/ksuid"
)

type UserAuth struct {
	AuthID         string            `json:"auth_id"`
	Email          string            `json:"email"`
	Role           []domain.UserRole `json:"role"`
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
	if record.User != nil {
		roles = record.User.RoleSystem
	}

	return UserAuth{
		AuthID:         record.ID,
		Email:          record.UserEmail,
		Role:           roles,
		RefreshTokenID: ksuid.New().String(),
	}
}

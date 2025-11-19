package domain

import "time"

type User struct {
	Email      string    `json:"email" bun:"email,pk"`
	RoleSystem string    `json:"role_system" bun:"role_system"`
	CreatedAt  time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bun:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" bun:"deleted_at"`
}

type Authentication struct {
	ID                string     `json:"id" bun:"id,pk"`
	UserEmail         string     `json:"user_email" bun:"user_email"`
	Password          string     `json:"password" bun:"password"`
	EncryptedPassword string     `json:"encrypted_password" bun:"encrypted_password"`
	RefreshTokenID    *string    `json:"refresh_token_id" bun:"refresh_token_id"`
	CreatedAt         time.Time  `json:"created_at" bun:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" bun:"updated_at"`
}
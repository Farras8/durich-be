package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleSales     UserRole = "sales"
	RoleWarehouse UserRole = "warehouse"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:user"`

	ID         string     `bun:",pk" json:"id"`
	Email      string     `bun:",unique,notnull" json:"email"`
	RoleSystem []UserRole `bun:",array" json:"role_system"`
	CreatedAt  time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt  *time.Time `bun:"," json:"deleted_at,omitempty"`
}

func (m *User) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.ID == "" {
			m.ID = ksuid.New().String()
		}
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

type Authentication struct {
	bun.BaseModel `bun:"table:authentications,alias:authentication"`

	ID                string    `bun:",pk" json:"id"`
	UserEmail         string    `bun:",nullzero" json:"user_email"`
	User              *User     `bun:"rel:belongs-to,join:user_email=email" json:"-"`
	Password          string    `bun:",nullzero" json:"-"`
	EncryptedPassword string    `bun:"," json:"-"`
	RefreshTokenID    *string   `bun:",nullzero" json:"-"`
	CreatedAt         time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

func (m *Authentication) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.ID == "" {
			m.ID = ksuid.New().String()
		}
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

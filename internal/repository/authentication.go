package repository

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/pkg/database"

	"github.com/uptrace/bun"
)

type AuthenticationRepository interface {
	Create(ctx context.Context, auth *domain.Authentication) error
	GetByEmail(ctx context.Context, email string) (domain.Authentication, error)
	GetByUserEmail(ctx context.Context, userEmail string) (domain.Authentication, error)
	GetByID(ctx context.Context, id string) (domain.Authentication, error)
	GetByRefreshTokenID(ctx context.Context, tokenID string) (domain.Authentication, error)
	Update(ctx context.Context, auth *domain.Authentication) error
}

type authenticationRepository struct {
	db *database.Database
}

func NewAuthenticationRepository(db *database.Database) AuthenticationRepository {
	return &authenticationRepository{db: db}
}

func (r *authenticationRepository) Create(ctx context.Context, auth *domain.Authentication) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(auth).Returning("id").Exec(ctx)
	return err
}

func (r *authenticationRepository) GetByEmail(ctx context.Context, email string) (domain.Authentication, error) {
	var data domain.Authentication
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&data).
		Column("id", "user_email", "password", "encrypted_password", "refresh_token_id", "created_at", "updated_at").
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "email", "role_system", "current_location_id")
		}).
		Where("user_email = ?", email).
		Scan(ctx)
	return data, err
}

func (r *authenticationRepository) GetByUserEmail(ctx context.Context, userEmail string) (domain.Authentication, error) {
	var data domain.Authentication
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&data).
		Column("id", "user_email", "password", "encrypted_password", "refresh_token_id", "created_at", "updated_at").
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "email", "role_system", "current_location_id")
		}).
		Where("user_email = ?", userEmail).
		Scan(ctx)
	return data, err
}

func (r *authenticationRepository) GetByID(ctx context.Context, id string) (domain.Authentication, error) {
	var data domain.Authentication
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&data).
		Column("id", "user_email", "password", "encrypted_password", "refresh_token_id", "created_at", "updated_at").
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "email", "role_system", "current_location_id")
		}).
		Where("authentication.id = ?", id).
		Scan(ctx)
	return data, err
}

func (r *authenticationRepository) GetByRefreshTokenID(ctx context.Context, tokenID string) (domain.Authentication, error) {
	var data domain.Authentication
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&data).
		Column("id", "user_email", "password", "encrypted_password", "refresh_token_id", "created_at", "updated_at").
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "email", "role_system", "current_location_id")
		}).
		Where("authentication.refresh_token_id = ?", tokenID).
		Scan(ctx)
	return data, err
}

func (r *authenticationRepository) Update(ctx context.Context, auth *domain.Authentication) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(auth).
		Where("id = ?", auth.ID).
		ExcludeColumn("created_at").
		Returning("id").
		Exec(ctx)
	return err
}

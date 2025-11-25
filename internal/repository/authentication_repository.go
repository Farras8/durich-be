package repository

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
)

type AuthenticationRepository interface {
	Create(ctx context.Context, auth *domain.Authentication) error
	GetByEmail(ctx context.Context, email string) (domain.Authentication, error)
	GetByUserEmail(ctx context.Context, userEmail string) (domain.Authentication, error)
	GetByID(ctx context.Context, id, ksuid *string) (domain.Authentication, error)
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
		Where("user_email = ?", email).
		Relation("User").
		Scan(ctx)
	return data, err
}

func (r *authenticationRepository) GetByUserEmail(ctx context.Context, userEmail string) (domain.Authentication, error) {
	var data domain.Authentication
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&data).
		Where("user_email = ?", userEmail).
		Relation("User").
		Scan(ctx)
	return data, err
}

func (r *authenticationRepository) GetByID(ctx context.Context, id, ksuid *string) (domain.Authentication, error) {
	var data domain.Authentication
	q := r.db.InitQuery(ctx).NewSelect().
		Model(&data).
		Relation("User")

	if id != nil {
		q.Where("authentication.id = ?", id)
	}

	if ksuid != nil {
		q.Where("authentication.refresh_token_id = ?", ksuid)
	}

	err := q.Scan(ctx)
	return data, err
}

func (r *authenticationRepository) Update(ctx context.Context, auth *domain.Authentication) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(auth).
		Where("id = ?", auth.ID).
		Returning("id").
		Exec(ctx)
	return err
}

package repository

import (
	"context"
	"durich-be/internal/domain"
	"github.com/uptrace/bun"
)

type AuthenticationRepository interface {
	Create(ctx context.Context, auth *domain.Authentication) error
	GetByUserEmail(ctx context.Context, userEmail string) (*domain.Authentication, error)
	Update(ctx context.Context, auth *domain.Authentication) error
}

type authenticationRepository struct {
	db *bun.DB
}

func NewAuthenticationRepository(db *bun.DB) AuthenticationRepository {
	return &authenticationRepository{db: db}
}

func (r *authenticationRepository) Create(ctx context.Context, auth *domain.Authentication) error {
	_, err := r.db.NewInsert().Model(auth).Exec(ctx)
	return err
}

func (r *authenticationRepository) GetByUserEmail(ctx context.Context, userEmail string) (*domain.Authentication, error) {
	auth := &domain.Authentication{}
	err := r.db.NewSelect().Model(auth).Where("user_email = ?", userEmail).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (r *authenticationRepository) Update(ctx context.Context, auth *domain.Authentication) error {
	_, err := r.db.NewUpdate().Model(auth).Where("user_email = ?", auth.UserEmail).Exec(ctx)
	return err
}
package repository

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByPhone(ctx context.Context, phone string) (domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&user).
		Where("email = ?", email).
		Scan(ctx)
	return user, err
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	var user domain.User
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model(user).
		Where("email = ?", user.Email).
		Exec(ctx)
	return err
}

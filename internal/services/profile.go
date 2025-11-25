package services

import (
	"context"
	"database/sql"

	"durich-be/internal/constants"
	"durich-be/internal/dto/requests"
	"durich-be/internal/repository"
	"durich-be/pkg/authentication"
	"durich-be/pkg/database"
	internal_err "durich-be/pkg/errors"

	"github.com/uptrace/bun"
)

type ProfileService interface {
	ChangePassword(ctx context.Context, email string, payload requests.ChangePasswordRequest) error
}

type profileService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthenticationRepository
}

func NewProfileService(
	userRepo repository.UserRepository,
	authRepo repository.AuthenticationRepository,
) ProfileService {
	return &profileService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (s *profileService) ChangePassword(ctx context.Context, email string, payload requests.ChangePasswordRequest) error {
	return database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return internal_err.NotFoundError(constants.UserNotFound)
		}

		authData, err := s.authRepo.GetByEmail(ctx, email)
		if err != nil {
			return internal_err.NotFoundError(constants.AuthDataNotFound)
		}

		isValid, err := authentication.VerifyPassword(payload.OldPassword, authData.Password)
		if err != nil {
			return err
		}

		if !isValid {
			return internal_err.ValidationError(constants.OldPasswordIncorrect)
		}

		hashedPassword, err := authentication.HashPassword(payload.NewPassword)
		if err != nil {
			return err
		}

		encryptedPassword, err := authentication.Encrypt(payload.NewPassword)
		if err != nil {
			return err
		}

		authData.Password = hashedPassword
		authData.EncryptedPassword = encryptedPassword

		err = s.authRepo.Update(ctx, &authData)
		if err != nil {
			return err
		}

		return nil
	})
}

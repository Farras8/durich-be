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

type MemberService interface {
	ResetUserPassword(ctx context.Context, payload requests.ResetPasswordRequest) error
}

type memberService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthenticationRepository
}

func NewMemberService(
	userRepo repository.UserRepository,
	authRepo repository.AuthenticationRepository,
) MemberService {
	return &memberService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (m *memberService) ResetUserPassword(ctx context.Context, payload requests.ResetPasswordRequest) error {
	return database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		user, err := m.userRepo.GetByEmail(ctx, payload.Email)
		if err != nil {
			return internal_err.NotFoundError(constants.UserNotFound)
		}

		authData, err := m.authRepo.GetByEmail(ctx, user.Email)
		if err != nil {
			return internal_err.NotFoundError(constants.AuthDataNotFound)
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

		err = m.authRepo.Update(ctx, &authData)
		if err != nil {
			return err
		}

		return nil
	})
}

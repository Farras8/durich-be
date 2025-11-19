package services

import (
	"context"
	"database/sql"
	"errors"

	"durich-be/internal/constants"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"durich-be/pkg/authentication"
	"durich-be/pkg/database"
	internal_err "durich-be/pkg/errors"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type AuthService interface {
	Register(ctx context.Context, payload requests.RegisterAdmin) (res response.RegisterResponse, err error)
	Login(ctx context.Context, payload requests.Login) (res response.LoginResponse, err error)
	Logout(ctx context.Context) (err error)
	RefreshToken(
		ctx context.Context,
		payload requests.RefreshToken,
	) (res response.RefreshTokenResponse, err error)
}

type authService struct {
	authRepo repository.AuthenticationRepository
	userRepo repository.UserRepository
}

func NewAuthService(authRepo repository.AuthenticationRepository, userRepo repository.UserRepository) AuthService {
	return &authService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (a *authService) Register(
	ctx context.Context,
	payload requests.RegisterAdmin,
) (res response.RegisterResponse, err error) {
	err = database.RunInTx(
		ctx,
		database.GetDB(),
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			existingUser, err := a.userRepo.GetByEmail(ctx, payload.Email)
			if err == nil && existingUser != nil && existingUser.Email != "" {
				return internal_err.ValidationError(constants.AuthEmailAlreadyExists)
			}

			hashedPassword, err := authentication.HashPassword(payload.Password)
			if err != nil {
				return err
			}

			encryptedPassword, err := authentication.Encrypt(payload.Password)
			if err != nil {
				return err
			}

			newUser := &domain.User{
				Email:      payload.Email,
				RoleSystem: "admin",
			}

			err = a.userRepo.Create(ctx, newUser)
			if err != nil {
				return err
			}

			newAuth := &domain.Authentication{
				UserEmail:         newUser.Email,
				Password:          hashedPassword,
				EncryptedPassword: encryptedPassword,
			}

			err = a.authRepo.Create(ctx, newAuth)
			if err != nil {
				return err
			}

			res = response.RegisterResponse{
				Email: newUser.Email,
			}

			return nil
		},
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (a *authService) Login(
	ctx context.Context,
	payload requests.Login,
) (res response.LoginResponse, err error) {
	err = database.RunInTx(
		ctx,
		database.GetDB(),
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			var authData domain.Authentication

			if payload.Email != "" {
				authData, err = a.authRepo.GetByEmail(ctx, payload.Email)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return internal_err.AuthError(constants.AuthPasswordInvalidOrEmailNotFound)
					}
					return err
				}
			} else if payload.Phone != "" {
				user, err := a.userRepo.GetByPhone(ctx, payload.Phone)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return internal_err.AuthError(constants.AuthPasswordInvalidOrEmailNotFound)
					}
					return err
				}

				authData, err = a.authRepo.GetByUserEmail(ctx, user.Email)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return internal_err.AuthError(constants.AuthPasswordInvalidOrEmailNotFound)
					}
					return err
				}
			} else {
				return internal_err.ValidationError("Email or phone number is required")
			}

			isValid, err := authentication.VerifyPassword(payload.Password, authData.Password)
			if err != nil {
				return err
			}

			if !isValid {
				return internal_err.AuthError(constants.AuthPasswordInvalidOrEmailNotFound)
			}

			tokenPayload := requests.ToTokenPayload(authData)
			pair, err := authentication.JWTAuth.GenerateTokenPair(tokenPayload, false)
			if err != nil {
				return err
			}

			authData.RefreshTokenID = &tokenPayload.RefreshTokenID
			err = a.authRepo.Update(ctx, &authData)
			if err != nil {
				return err
			}

			res = response.LoginResponse{
				AccessToken:  pair.AccessToken,
				RefreshToken: pair.RefreshToken,
				Role:         tokenPayload.Role,
			}

			return nil
		},
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (a *authService) Logout(ctx context.Context) (err error) {
	err = database.RunInTx(
		ctx,
		database.GetDB(),
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			authToken := authentication.GetUserDataFromToken(ctx)
			if authToken.AuthID == "" {
				return internal_err.AuthError(constants.DataNotFound)
			}

			authData, err := a.authRepo.GetByID(ctx, &authToken.AuthID, nil)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return internal_err.AuthError(constants.DataNotFound)
				}
				return err
			}

			authData.RefreshTokenID = nil
			err = a.authRepo.Update(ctx, &authData)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) RefreshToken(
	ctx context.Context,
	payload requests.RefreshToken,
) (res response.RefreshTokenResponse, err error) {
	claims, err := authentication.JWTAuth.VerifyRefreshToken(payload.RefreshToken)
	if err != nil {
		return res, err
	}
	if claims == nil {
		return res, internal_err.AuthError(constants.AuthInvalidToken)
	}

	auth, err := a.authRepo.GetByID(ctx, nil, &claims.TokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, internal_err.AuthError(constants.DataNotFound)
		}
		return res, err
	}

	tokenPayload := requests.ToTokenPayload(auth)
	tokenPayload.RefreshTokenID = ksuid.New().String()

	claimsRefresh, err := authentication.JWTAuth.GenerateTokenPair(tokenPayload, true)
	if err != nil {
		return res, err
	}

	return response.RefreshTokenResponse{
		AccessToken: claimsRefresh.AccessToken,
	}, err
}

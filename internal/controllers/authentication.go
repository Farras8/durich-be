package controllers

import (
	"net/http"

	"durich-be/internal/constants"
	"durich-be/internal/dto/requests"
	"durich-be/internal/services"
	"durich-be/pkg/errors"
	"durich-be/pkg/http/response"
	"durich-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authSrv services.AuthService) AuthController {
	return AuthController{
		authService: authSrv,
	}
}

func (ctl *AuthController) Register(ctx *gin.Context) {
	var payload requests.RegisterAdmin
	if err := utils.BindData(ctx, &payload); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.authService.Register(ctx, payload)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusCreated, constants.AuthRegisterSuccess, res)
}

func (ctl *AuthController) RegisterWarehouse(ctx *gin.Context) {
	var payload requests.RegisterAdmin
	if err := utils.BindData(ctx, &payload); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.authService.RegisterWarehouse(ctx, payload)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusCreated, constants.AuthRegisterSuccess, res)
}

func (ctl *AuthController) RegisterSales(ctx *gin.Context) {
	var payload requests.RegisterAdmin
	if err := utils.BindData(ctx, &payload); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.authService.RegisterSales(ctx, payload)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusCreated, constants.AuthRegisterSuccess, res)
}

func (ctl *AuthController) Login(ctx *gin.Context) {
	var auth requests.Login
	if err := utils.BindData(ctx, &auth); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.authService.Login(ctx, auth)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, constants.AuthLoginSuccess, res)
}

func (ctl *AuthController) Logout(ctx *gin.Context) {
	err := ctl.authService.Logout(ctx)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, constants.AuthLogoutSuccess, nil)
}

func (ctl *AuthController) RefreshToken(ctx *gin.Context) {
	var auth requests.RefreshToken
	if err := utils.BindData(ctx, &auth); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.authService.RefreshToken(ctx, auth)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, constants.AuthRefreshTokenSuccess, res)
}
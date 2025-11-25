package controllers

import (
	"net/http"

	"durich-be/internal/constants"
	"durich-be/internal/dto/requests"
	"durich-be/internal/services"
	"durich-be/pkg/authentication"
	"durich-be/pkg/errors"
	http_response "durich-be/pkg/http/response"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	profileService services.ProfileService
}

func NewProfileController(profileService services.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

func (ctrl *ProfileController) ChangePassword(ctx *gin.Context) {
	userAuth := authentication.GetUserDataFromToken(ctx)
	if userAuth.Email == "" {
		http_response.SendError(ctx, errors.AuthError("invalid token"))
		return
	}

	var payload requests.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := ctrl.profileService.ChangePassword(ctx.Request.Context(), userAuth.Email, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, constants.PasswordChangedSuccess, nil)
}

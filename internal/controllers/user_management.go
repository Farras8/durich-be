package controllers

import (
	"net/http"

	"durich-be/internal/constants"
	"durich-be/internal/dto/requests"
	"durich-be/internal/services"
	"durich-be/pkg/errors"
	http_response "durich-be/pkg/http/response"
	"durich-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MemberController struct {
	memberService services.MemberService
}

func NewMemberController(memberService services.MemberService) *MemberController {
	return &MemberController{
		memberService: memberService,
	}
}

func (ctrl *MemberController) ResetPassword(ctx *gin.Context) {
	var payload requests.ResetPasswordRequest
	if err := utils.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := ctrl.memberService.ResetUserPassword(ctx.Request.Context(), payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, constants.PasswordResetSuccess, nil)
}

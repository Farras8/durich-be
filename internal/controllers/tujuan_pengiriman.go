package controllers

import (
	"net/http"

	"durich-be/internal/dto/requests"
	"durich-be/internal/services"
	"durich-be/pkg/authentication"
	"durich-be/pkg/errors"
	"durich-be/pkg/http/response"
	"durich-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TujuanPengirimanController struct {
	service services.TujuanPengirimanService
}

func NewTujuanPengirimanController(service services.TujuanPengirimanService) TujuanPengirimanController {
	return TujuanPengirimanController{service: service}
}

func (c *TujuanPengirimanController) Create(ctx *gin.Context) {
	var req requests.CreateTujuanPengirimanRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	result, err := c.service.Create(ctx.Request.Context(), req, locationID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusCreated, "Tujuan pengiriman created successfully", result)
}

func (c *TujuanPengirimanController) GetAll(ctx *gin.Context) {
	result, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Tujuan pengiriman retrieved successfully", result)
}

func (c *TujuanPengirimanController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Tujuan pengiriman retrieved successfully", result)
}

func (c *TujuanPengirimanController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.UpdateTujuanPengirimanRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	result, err := c.service.Update(ctx.Request.Context(), id, req, locationID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Tujuan pengiriman updated successfully", result)
}

func (c *TujuanPengirimanController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	err := c.service.Delete(ctx.Request.Context(), id, locationID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Tujuan pengiriman deleted successfully", nil)
}

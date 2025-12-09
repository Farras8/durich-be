package controllers

import (
	"durich-be/internal/dto/requests"
	"durich-be/internal/services"
	"durich-be/pkg/authentication"
	"durich-be/pkg/errors"
	"durich-be/pkg/http/response"
	"durich-be/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SalesController struct {
	service services.SalesService
}

func NewSalesController(service services.SalesService) *SalesController {
	return &SalesController{service: service}
}

func (c *SalesController) Create(ctx *gin.Context) {
	var req requests.SalesCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusCreated, "Sales invoice created successfully", res)
}

func (c *SalesController) GetList(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	tipeJual := ctx.Query("tipe_jual")

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	res, err := c.service.GetList(ctx.Request.Context(), startDate, endDate, tipeJual, locationID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Sales list retrieved successfully", res)
}

func (c *SalesController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Sales detail retrieved successfully", res)
}

func (c *SalesController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.SalesUpdateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := c.service.Update(ctx.Request.Context(), id, req); err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Sales info updated successfully", nil)
}

func (c *SalesController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID
	
	// Get user role from auth token
	userRole := ""
	if len(userAuth.Role) > 0 {
		userRole = string(userAuth.Role[0])
	}

	if err := c.service.Delete(ctx.Request.Context(), id, locationID, userRole); err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Sales voided successfully", nil)
}

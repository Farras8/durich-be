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

type ShipmentController struct {
	service services.ShipmentService
}

func NewShipmentController(service services.ShipmentService) *ShipmentController {
	return &ShipmentController{service: service}
}

func (c *ShipmentController) Create(ctx *gin.Context) {
	var req requests.ShipmentCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	if userAuth.UserID == "" {
		response.SendError(ctx, errors.AuthError("Invalid token: missing user_id. Please login again to refresh your token."))
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), req, userAuth.UserID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusCreated, "Shipment created successfully", res)
}

func (c *ShipmentController) GetList(ctx *gin.Context) {
	tujuan := ctx.Query("tujuan")
	status := ctx.Query("status")

	res, err := c.service.GetList(ctx.Request.Context(), tujuan, status)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Shipments retrieved successfully", res)
}

func (c *ShipmentController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Shipment detail retrieved successfully", res)
}

func (c *ShipmentController) AddItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.ShipmentAddItemRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := c.service.AddItem(ctx.Request.Context(), id, req); err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Item added successfully", nil)
}

func (c *ShipmentController) RemoveItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.ShipmentRemoveItemRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := c.service.RemoveItem(ctx.Request.Context(), id, req.DetailID); err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Item removed successfully", nil)
}

func (c *ShipmentController) Finalize(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Finalize(ctx.Request.Context(), id); err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Shipment finalized successfully", nil)
}

package controllers

import (
	"durich-be/internal/services"
	"durich-be/pkg/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TraceabilityController struct {
	service services.TraceabilityService
}

func NewTraceabilityController(service services.TraceabilityService) *TraceabilityController {
	return &TraceabilityController{service: service}
}

func (c *TraceabilityController) TraceLot(ctx *gin.Context) {
	lotID := ctx.Param("id")
	res, err := c.service.TraceLot(ctx.Request.Context(), lotID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Lot traceability retrieved successfully", res)
}

func (c *TraceabilityController) TraceFruit(ctx *gin.Context) {
	fruitID := ctx.Param("buah_raw_id")
	res, err := c.service.TraceFruit(ctx.Request.Context(), fruitID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Fruit traceability retrieved successfully", res)
}

func (c *TraceabilityController) TraceShipment(ctx *gin.Context) {
	shipmentID := ctx.Param("id")
	res, err := c.service.TraceShipment(ctx.Request.Context(), shipmentID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Shipment traceability retrieved successfully", res)
}

package controllers

import (
	"durich-be/internal/services"
	"durich-be/pkg/http/response"
	"net/http"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service services.DashboardService
}

func NewDashboardController(service services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

func (c *DashboardController) GetStokDashboard(ctx *gin.Context) {
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	res, err := c.service.GetStokDashboard(ctx.Request.Context(), dateFrom, dateTo)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Stok dashboard retrieved successfully", res)
}

func (c *DashboardController) GetSalesDashboard(ctx *gin.Context) {
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	res, err := c.service.GetSalesDashboard(ctx.Request.Context(), dateFrom, dateTo)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Sales dashboard retrieved successfully", res)
}

func (c *DashboardController) GetWarehouseData(ctx *gin.Context) {
	res, err := c.service.GetWarehouseData(ctx.Request.Context())
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Warehouse data retrieved successfully", res)
}

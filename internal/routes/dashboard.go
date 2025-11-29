package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterDashboard(router *gin.RouterGroup, ctl *controllers.DashboardController) {
	stokGroup := router.Group("/dashboard/stok")
	stokGroup.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse))
	{
		stokGroup.GET("", ctl.GetStokDashboard)
	}

	salesGroup := router.Group("/dashboard/sales")
	salesGroup.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin, domain.RoleSales))
	{
		salesGroup.GET("", ctl.GetSalesDashboard)
	}

	warehouseDataGroup := router.Group("/dashboard/warehouse-data")
	warehouseDataGroup.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse))
	{
		warehouseDataGroup.GET("", ctl.GetWarehouseData)
	}
}

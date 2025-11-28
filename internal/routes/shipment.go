package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterShipment(router *gin.RouterGroup, ctl *controllers.ShipmentController) {
	group := router.Group("/shipments")
	group.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse))
	{
		group.POST("", ctl.Create)
		group.GET("", ctl.GetList)
		group.GET("/:id", ctl.GetByID)
		group.POST("/:id/items", ctl.AddItem)
		group.DELETE("/:id/items", ctl.RemoveItem)
		group.POST("/:id/finalize", ctl.Finalize)
	}

	salesGroup := router.Group("/shipments")
	salesGroup.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin, domain.RoleSales))
	{
		salesGroup.PATCH("/:id/status", ctl.UpdateStatus)
	}
}

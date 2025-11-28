package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterTraceability(router *gin.RouterGroup, ctl *controllers.TraceabilityController) {
	group := router.Group("/trace")
	group.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse, domain.RoleSales))
	{
		group.GET("/lot/:id", ctl.TraceLot)
		group.GET("/fruit/:buah_raw_id", ctl.TraceFruit)
		group.GET("/shipment/:id", ctl.TraceShipment)
	}
}

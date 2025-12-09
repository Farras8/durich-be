package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterTujuanPengiriman(router *gin.RouterGroup, ctl controllers.TujuanPengirimanController) {
	group := router.Group("/tujuan-pengiriman")
	group.Use(middlewares.TokenAuthMiddleware())
	{
		// Read Access: Admin, Warehouse, Sales
		group.GET("", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse, domain.RoleSales), ctl.GetAll)
		group.GET("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse, domain.RoleSales), ctl.GetByID)

		// Write Access: Admin & Warehouse (But logic restricts to Central User only)
		group.POST("", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.Create)
		group.PUT("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.Update)
		group.DELETE("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.Delete)
	}
}

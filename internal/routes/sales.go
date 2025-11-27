package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterSales(router *gin.RouterGroup, ctl *controllers.SalesController) {
	group := router.Group("/sales")
	group.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin))
	{
		group.POST("", ctl.Create)
		group.GET("", ctl.GetList)
		group.GET("/:id", ctl.GetByID)
		group.PUT("/:id", ctl.Update)
		group.DELETE("/:id", ctl.Delete)
	}
}

package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterLotRoutes(r *gin.RouterGroup, lotController *controllers.LotController) {
	lots := r.Group("/lots")
	lots.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse))
	{
		lots.POST("", lotController.Create)
		lots.GET("", lotController.GetList)
		lots.GET("/:id", lotController.GetDetail)
		lots.POST("/:id/items", lotController.AddItems)
		lots.DELETE("/:id/items", lotController.RemoveItem)
		lots.POST("/:id/finalize", lotController.Finalize)
	}
}

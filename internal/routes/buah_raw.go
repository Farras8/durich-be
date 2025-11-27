package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterBuahRaw(router *gin.RouterGroup, ctl controllers.BuahRawController) {
	group := router.Group("/buah-raw")
	group.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse))
	{
		group.POST("", ctl.Create)
		group.POST("/bulk", ctl.BulkCreate)
		group.GET("", ctl.GetList)
		group.GET("/:id", ctl.GetDetail)
		group.PUT("/:id", ctl.Update)
		group.DELETE("/:id", ctl.Delete)
	}
}

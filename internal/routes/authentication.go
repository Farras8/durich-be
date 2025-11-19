package routes

import (
	"durich-be/internal/controllers"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(router *gin.RouterGroup, authCtl controllers.AuthController) {
	admin := router.Group("/admin")
	{
		admin.POST("/register-admin", authCtl.Register)
	}

	auth := router.Group("/authentications")
	{
		auth.POST("/login", authCtl.Login)
		auth.POST("/refresh-token", authCtl.RefreshToken)

		auth.Use(middlewares.TokenAuthMiddleware())
		auth.POST("/logout", authCtl.Logout)
	}
}

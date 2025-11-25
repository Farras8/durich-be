package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(router *gin.RouterGroup, authCtl controllers.AuthController, profileCtl *controllers.ProfileController, memberCtl *controllers.MemberController) {
	admin := router.Group("/admin")
	{
		admin.POST("/register-admin", authCtl.Register)
		admin.POST("/register-warehouse", authCtl.RegisterWarehouse)
		admin.POST("/register-sales", authCtl.RegisterSales)

		admin.Use(middlewares.TokenAuthMiddleware(), middlewares.RoleHandler(domain.RoleAdmin))
		admin.POST("/users/reset-password", memberCtl.ResetPassword)
	}

	auth := router.Group("/authentications")
	{
		auth.POST("/login", authCtl.Login)
		auth.POST("/refresh-token", authCtl.RefreshToken)

		auth.Use(middlewares.TokenAuthMiddleware())
		auth.POST("/logout", authCtl.Logout)
	}

	profile := router.Group("/profile")
	{
		profile.Use(middlewares.TokenAuthMiddleware())
		profile.PUT("/password", profileCtl.ChangePassword)
	}
}

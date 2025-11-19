package main

import (
	"log"
	"durich-be/internal/controllers"
	"durich-be/internal/repository"
	"durich-be/internal/routes"
	"durich-be/internal/services"
	"durich-be/pkg/authentication"
	"durich-be/pkg/config"
	"durich-be/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := database.Connect(cfg.Database)

	jwtService := authentication.NewJWTService(
		cfg.Authentication.AccessSecretKey,
		cfg.Authentication.RefreshSecretKey,
		cfg.Authentication.Issuer,
		cfg.Authentication.AccessTokenExpiry,
		cfg.Authentication.RefreshTokenExpiry,
	)

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthenticationRepository(db)

	authService := services.NewAuthService(userRepo, authRepo, jwtService)

	authController := controllers.NewAuthController(authService)

	router := gin.Default()

	routes.SetupAuthRoutes(router, authController)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	router.Run(":" + cfg.Server.Port)
}
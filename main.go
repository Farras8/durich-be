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

	database.InitDB(cfg.Database)
	db := database.GetDB()

	authentication.SetupKey(cfg.Authentication.EncryptKey)
	authentication.NewJWTManager(authentication.JWTOptions{
		AccessSecret:       cfg.Authentication.AccessSecretKey,
		RefreshSecret:      cfg.Authentication.RefreshSecretKey,
		Issuer:             cfg.Authentication.Issuer,
		ExpiryAccessToken:  cfg.Authentication.AccessTokenExpiry,
		ExpiryRefreshToken: cfg.Authentication.RefreshTokenExpiry,
	})

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthenticationRepository(db)

	authService := services.NewAuthService(userRepo, authRepo)
	profileService := services.NewProfileService(userRepo, authRepo)
	memberService := services.NewMemberService(userRepo, authRepo)

	authController := controllers.NewAuthController(authService)
	profileController := controllers.NewProfileController(profileService)
	memberController := controllers.NewMemberController(memberService)

	router := gin.Default()

	v1 := router.Group("/v1")
	routes.RegisterAuth(v1, authController, profileController, memberController)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	router.Run(":" + cfg.Server.Port)
}

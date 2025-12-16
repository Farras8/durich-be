package main

import (
	"log"
	"strings"
	"time"

	"durich-be/internal/controllers"
	"durich-be/internal/repository"
	"durich-be/internal/routes"
	"durich-be/internal/services"
	"durich-be/pkg/authentication"
	"durich-be/pkg/config"
	"durich-be/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, _ := config.LoadConfig()

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
	buahRawRepo := repository.NewBuahRawRepository(db)
	masterDataRepo := repository.NewMasterDataRepository(db)
	lotRepo := repository.NewLotRepository(db.DB)
	shipmentRepo := repository.NewShipmentRepository()
	tujuanPengirimanRepo := repository.NewTujuanPengirimanRepository(db)
	salesRepo := repository.NewSalesRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	traceabilityRepo := repository.NewTraceabilityRepository(db)

	authService := services.NewAuthService(userRepo, authRepo)
	profileService := services.NewProfileService(userRepo, authRepo)
	memberService := services.NewMemberService(userRepo, authRepo)
	buahRawService := services.NewBuahRawService(buahRawRepo)
	masterDataService := services.NewMasterDataService(masterDataRepo)
	lotService := services.NewLotService(db.DB, lotRepo, buahRawRepo)
	shipmentService := services.NewShipmentService(db, shipmentRepo, tujuanPengirimanRepo)
	tujuanPengirimanService := services.NewTujuanPengirimanService(tujuanPengirimanRepo)
	salesService := services.NewSalesService(salesRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)
	traceabilityService := services.NewTraceabilityService(traceabilityRepo)

	authController := controllers.NewAuthController(authService)
	profileController := controllers.NewProfileController(profileService)
	memberController := controllers.NewMemberController(memberService)
	buahRawController := controllers.NewBuahRawController(buahRawService)
	masterDataController := controllers.NewMasterDataController(masterDataService)
	lotController := controllers.NewLotController(lotService)
	shipmentController := controllers.NewShipmentController(shipmentService)
	tujuanPengirimanController := controllers.NewTujuanPengirimanController(tujuanPengirimanService)
	salesController := controllers.NewSalesController(salesService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	traceabilityController := controllers.NewTraceabilityController(traceabilityService)

	router := gin.Default()

	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Allow all localhost and 127.0.0.1 origins with any port
			return strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := router.Group("/v1")
	routes.RegisterAuth(v1, authController, profileController, memberController)
	routes.RegisterBuahRaw(v1, buahRawController)
	routes.RegisterMasterData(v1, masterDataController)
	routes.RegisterLotRoutes(v1, lotController)
	routes.RegisterShipment(v1, shipmentController)
	routes.RegisterTujuanPengiriman(v1, tujuanPengirimanController)
	routes.RegisterSales(v1, salesController)
	routes.RegisterDashboard(v1, dashboardController)
	routes.RegisterTraceability(v1, traceabilityController)

	log.Printf("Server running on port %s", cfg.Server.Port)
	log.Fatal(router.Run(":" + cfg.Server.Port))
}

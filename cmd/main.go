package main

import (
	"log"

	"gotickets/internal/config"
	"gotickets/internal/handler"
	"gotickets/internal/middlewares"
	"gotickets/internal/models"
	"gotickets/internal/repository"
	"gotickets/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// 1. Load config
	cfg := config.LoadEnv()

	// 2. Connect Database
	db := config.ConnectDatabase(cfg)

	// AutoMigrate the schema
	err := db.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	// 4. Initialize Services
	authService := service.NewAuthService(userRepo, cfg)
	zoneService := service.NewZoneService(zoneRepo, reservationRepo)
	reservationService := service.NewReservationService(reservationRepo)

	// 5. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)

	// 6. Setup Echo
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// 7. Routes
	api := e.Group("/api/v1")

	// Auth routes
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// Zones routes
	zoneGroup := api.Group("/zones")
	zoneGroup.GET("", zoneHandler.GetAllZones)
	zoneGroup.GET("/:id", zoneHandler.GetZoneByID)
	
	// Admin only zone routes
	adminZoneGroup := zoneGroup.Group("", middlewares.JWTMiddleware(cfg.JwtSecret), middlewares.RoleMiddleware("admin"))
	adminZoneGroup.POST("", zoneHandler.CreateZone)

	// Reservations routes
	reservationGroup := api.Group("/reservations", middlewares.JWTMiddleware(cfg.JwtSecret))
	
	// Authenticated users (driver, admin)
	reservationGroup.POST("", reservationHandler.CreateReservation, middlewares.RoleMiddleware("driver", "admin"))
	reservationGroup.GET("/my-reservations", reservationHandler.GetMyReservations, middlewares.RoleMiddleware("driver", "admin"))
	reservationGroup.DELETE("/:id", reservationHandler.CancelReservation, middlewares.RoleMiddleware("driver", "admin"))

	// Admin only reservation routes
	adminReservationGroup := reservationGroup.Group("", middlewares.RoleMiddleware("admin"))
	adminReservationGroup.GET("", reservationHandler.GetAllReservations)

	// 8. Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}

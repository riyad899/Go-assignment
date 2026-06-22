package server

import (
	"gotickets/internal/config"
	"gotickets/internal/domain/user"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate(&user.User{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "running")
	})

	//routes
	user.RegisterRoutes(e, db, cfg)

	port := os.Getenv("PORT")

	if port == "" {
		port = cfg.Port

		if port == "" {
			port = "8080"
		}
	}

	log.Printf("Starting server on port %s", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

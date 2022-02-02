package app

import (
	"fmt"
	"strings"

	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/internal/rest"
	"github.com/adrianolmedo/go-restapi/internal/service"
	"github.com/adrianolmedo/go-restapi/internal/storage"
	"github.com/adrianolmedo/go-restapi/jwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(cfg *config.Config) error {
	// Load authentication credentials.
	err := jwt.LoadFiles("app.sra", "app.sra.pub")
	if err != nil {
		return fmt.Errorf("certificates could not be loaded: %v", err)
	}

	// Echo framework.
	e := echo.New()

	// - Load Echo middlewares.
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// - CORS restricted with GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(cfg.CORS, ","),
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Prepare services.
	s := storage.New(cfg.Database)

	svc, err := service.New(s)
	if err != nil {
		return err
	}

	// - Call routes.
	rest.Routes(e, *svc)

	// - Up server.
	return e.Start(":" + cfg.Port)
}

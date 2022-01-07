package app

import (
	"log"
	"strings"

	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/internal/rest"
	"github.com/adrianolmedo/go-restapi-practice/internal/service"
	"github.com/adrianolmedo/go-restapi-practice/jwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(cfgPath string) {
	// Load config file.
	cfg, err := config.Init(cfgPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Load authentication credentials.
	err = jwt.LoadFiles("app.sra", "app.sra.pub")
	if err != nil {
		log.Fatalf("Certificates could not be loaded: %v", err)
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
	svc, err := service.New(cfg.Database)
	if err != nil {
		log.Printf("%v\n", err)
	}

	// - Call routes.
	rest.Routes(e, *svc)

	// - Up server.
	err = e.Start(cfg.Address())
	if err != nil {
		log.Printf("Error server: %v\n", err)
	}
}

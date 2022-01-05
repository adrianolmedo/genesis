package app

import (
	"log"
	"strings"

	"github.com/adrianolmedo/go-restapi-practice/config"
	"github.com/adrianolmedo/go-restapi-practice/infra/jwt"
	"github.com/adrianolmedo/go-restapi-practice/internal/server/rest"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(cfgPath string) {
	var repos *storage.Repositories

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

	// - Set up storage from configuration loaded.
	repos, err = storage.NewRepositories(cfg.Database)
	if err != nil {
		log.Fatalf("Error from storage: %v\n", err)
	}

	// Echo framework.
	e := echo.New()

	// - Load middlewares.
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// - CORS restricted with GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(cfg.CORS, ","),
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// - Call routes.
	rest.Routes(e, repos)

	// - Up server.
	err = e.Start(cfg.Address())
	if err != nil {
		log.Printf("Error server: %v\n", err)
	}
}

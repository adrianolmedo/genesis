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
	"github.com/labstack/gommon/log"
)

func Run(cfg *config.Config) error {
	// Load authentication credentials.
	err := jwt.LoadFiles("app.rsa", "app.rsa.pub")
	if err != nil {
		return fmt.Errorf("certificates could not be loaded: %v", err)
	}

	// Echo framework.
	e := echo.New()

	// - Load Echo middlewares.
	e.Use(middleware.Recover())

	// - Echo logger.
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}))

	// - CORS restricted with GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(cfg.CORS, ","),
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Prepare services.
	svc, err := service.New(storage.New(cfg.Database))
	if err != nil {
		return err
	}

	// - Call routes.
	rest.Routes(e, *svc)

	// - Up server.
	return e.Start(":" + cfg.Port)
}

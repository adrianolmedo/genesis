package main

import (
	"fmt"

	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/delivery"
	"github.com/adrianolmedo/go-restapi/postgres"

	"github.com/gofiber/fiber/v2"
)

func run(cfg *config.Config) error {
	app := fiber.New()

	// Prepare storage.
	stg, err := postgres.NewStorage(cfg.DB)
	if err != nil {
		return fmt.Errorf("error from storage: %v", err)
	}

	delivery.Routes(app, stg)

	// Up server.
	return app.Listen(":" + cfg.Port)
}

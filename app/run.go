package app

import (
	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/product"
	"github.com/adrianolmedo/go-restapi/user"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg *config.Config) error {
	app := fiber.New()

	// Prepare storages.
	userRepo, err := user.NewRepository(cfg.Database)
	if err != nil {
		return err
	}

	// Prepare services.
	svc := user.NewService(userRepo)

	// Call routes.
	user.Routes(app, *svc)
	product.Routes(app)

	// Up server.
	return app.Listen(":" + cfg.Port)
}

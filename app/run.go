package app

import (
	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/storage"
	"github.com/adrianolmedo/go-restapi/store"
	"github.com/adrianolmedo/go-restapi/user"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg *config.Config) error {
	app := fiber.New()

	// Prepare storage.
	db, err := storage.New(cfg.DB)
	if err != nil {
		return err
	}

	// Prepare services.
	userSvc := user.NewService(db)
	storeSvc := store.NewService(db)

	// Call routes.
	user.Routes(app, *userSvc)
	store.Routes(app, *storeSvc)

	// Up server.
	return app.Listen(":" + cfg.Port)
}

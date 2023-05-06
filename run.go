package main

import (
	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/storage"
	"github.com/adrianolmedo/go-restapi/store"
	"github.com/adrianolmedo/go-restapi/user"

	"github.com/gofiber/fiber/v2"
)

func run(cfg *config.Config) error {
	app := fiber.New()

	// Prepare storage.
	db, err := storage.New(cfg.DB)
	if err != nil {
		return err
	}

	// Prepare repositories.
	userRepo := storage.NewUserRepository(db)
	storeRepo := store.NewRepository(db)

	// Prepare services.
	userSvc := user.NewService(userRepo)
	storeSvc := store.NewService(storeRepo)

	// Call routes.
	user.Routes(app, userSvc)
	store.Routes(app, storeSvc)

	// Up server.
	return app.Listen(":" + cfg.Port)
}

package rest

import (
	"github.com/adrianolmedo/go-restapi/billing"
	"github.com/adrianolmedo/go-restapi/postgres"
	"github.com/adrianolmedo/go-restapi/rest/middleware"
	"github.com/adrianolmedo/go-restapi/store"
	"github.com/adrianolmedo/go-restapi/user"

	"github.com/gofiber/fiber/v2"
)

type services struct {
	User    user.Service
	Store   store.Service
	Billing billing.Service
}

// Routes call all routes with its handlers.
func Routes(strg *postgres.Storage) *fiber.App {
	s := &services{
		User:    user.NewService(strg.User),
		Store:   store.NewService(strg.Product),
		Billing: billing.NewService(strg.Invoice),
	}

	app := fiber.New()

	g := app.Group("/v1/users")
	g.Use(middleware.Auth)
	g.Get("", listUsers(s))
	g.Put("/:id", updateUser(s))
	g.Delete("/:id", deleteUser(s))

	app.Post("v1/login", loginUser(s))
	app.Post("v1/users", signUpUser(s))
	app.Get("/v1/users/:id", findUser(s))

	app.Get("/v1/products", listProducts(s))
	app.Get("/v1/products/:id", findProduct(s))

	app.Post("/v1/products", addProduct(s), middleware.Auth)
	app.Put("/v1/products/:id", updateProduct(s), middleware.Auth)
	app.Delete("/v1/products/:id", deleteProduct(s), middleware.Auth)

	app.Post("v1/invoices", generateInvoice(s), middleware.Auth)
	return app
}

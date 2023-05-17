package rest

import (
	"github.com/adrianolmedo/go-restapi/billing"
	"github.com/adrianolmedo/go-restapi/postgres"
	"github.com/adrianolmedo/go-restapi/store"
	"github.com/adrianolmedo/go-restapi/user"

	"github.com/gofiber/fiber/v2"
)

type services struct {
	User    user.Service
	Store   store.Service
	Billing billing.Service
}

func Routes(strg *postgres.Storage) *fiber.App {
	app := fiber.New()

	s := &services{
		User:    user.NewService(strg.User),
		Store:   store.NewService(strg.Product),
		Billing: billing.NewService(strg.Invoice),
	}

	app.Post("v1/users", signUpUser(s))
	app.Get("v1/users", listUsers(s))
	app.Get("/v1/users/:id", findUser(s))
	app.Put("v1/users/:id", updateUser(s))
	app.Delete("v1/users/:id", deleteUser(s))

	app.Post("/v1/products", addProduct(s))
	app.Get("/v1/products", listProducts(s))
	app.Get("/v1/products/:id", findProduct(s))
	app.Put("/v1/products/:id", updateProduct(s))
	app.Delete("/v1/products/:id", deleteProduct(s))

	app.Post("v1/invoices", generateInvoice(s))
	return app
}

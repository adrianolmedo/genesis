package rest

import (
	app "github.com/adrianolmedo/go-restapi"
	"github.com/adrianolmedo/go-restapi/postgres"
	"github.com/adrianolmedo/go-restapi/rest/middleware"

	"github.com/gofiber/fiber/v2"
)

// Routes call all routes with its handlers.
func Routes(strg *postgres.Storage) *fiber.App {
	s := app.NewServices(strg)
	f := fiber.New()

	g := f.Group("/v1/users")
	g.Use(middleware.Auth)
	g.Get("", listUsers(s))
	g.Put("/:id", updateUser(s))
	g.Delete("/:id", deleteUser(s))

	f.Post("v1/login", loginUser(s))
	f.Post("v1/users", signUpUser(s))
	f.Get("/v1/users/:id", findUser(s))

	f.Get("/v1/products", listProducts(s))
	f.Get("/v1/products/:id", findProduct(s))

	f.Post("/v1/products", addProduct(s), middleware.Auth)
	f.Put("/v1/products/:id", updateProduct(s), middleware.Auth)
	f.Delete("/v1/products/:id", deleteProduct(s), middleware.Auth)

	f.Post("v1/invoices", generateInvoice(s), middleware.Auth)
	return f
}

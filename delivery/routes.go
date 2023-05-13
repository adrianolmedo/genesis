package delivery

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

func Routes(f *fiber.App, stg *postgres.Storage) {
	s := &services{
		User:    user.NewService(stg.User),
		Store:   store.NewService(stg.Product),
		Billing: billing.NewService(stg.Invoice),
	}

	f.Post("v1/users", signUpUser(s))
	f.Get("v1/users", listUsers(s))
	f.Get("/v1/users/:id", findUser(s))
	f.Put("v1/users/:id", updateUser(s))
	f.Delete("v1/users/:id", deleteUser(s))

	f.Post("/v1/products", addProduct(s))
	f.Get("/v1/products", listProducts(s))
	f.Get("/v1/products/:id", findProduct(s))
	f.Put("/v1/products/:id", updateProduct(s))
	f.Delete("/v1/products/:id", deleteProduct(s))
}

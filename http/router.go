package http

import (
	"net/http"

	"github.com/adrianolmedo/genesis/app"
	_ "github.com/adrianolmedo/genesis/docs"
	"github.com/adrianolmedo/genesis/http/jwt"
	"github.com/adrianolmedo/genesis/postgres"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

//	@title			Genesis REST API
//	@version		1.0
//	@description	This is a sample server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Adrián Olmedo
//	@contact.url	https://twitter.com/adrianolmedo
//	@contact.email	adrianolmedo.ve@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:3000
// @BasePath	/v1/
func Router(strg *postgres.Storage) *fiber.App {
	s := app.NewServices(strg)
	f := fiber.New()

	f.Post("/v1/login", loginUser(s))
	f.Post("/v1/users", signUpUser(s))
	f.Get("/v1/users/:id", findUser(s))

	f.Get("/v1/users", authWare, listUsers(s))
	f.Put("/v1/users/:id", authWare, updateUser(s))
	f.Delete("/v1/users/:id", authWare, deleteUser(s))

	f.Post("/v1/customers", createCustomer(s))
	f.Get("/v1/customers", listCustomers(s))
	f.Delete("v1/customers/:id", deleteCustomer(s))

	f.Get("/v1/products", listProducts(s))
	f.Get("/v1/products/:id", findProduct(s))

	f.Post("/v1/products", authWare, addProduct(s))
	f.Put("/v1/products/:id", authWare, updateProduct(s))
	f.Delete("/v1/products/:id", authWare, deleteProduct(s))

	f.Post("/v1/invoices", authWare, generateInvoice(s))

	f.Get("/swagger/*", swagger.WrapHandler)

	return f
}

// authWare middleware for handlers that require user login.
func authWare(c *fiber.Ctx) error {
	token := c.Request().Header.Peek("Authorization")
	_, err := jwt.Verify(string(token))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(respError{Msg: "you don't have authorization"})
	}
	return c.Next()
}

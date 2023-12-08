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

	f.Get("/v1/users", authMiddleware, listUsers(s))
	f.Put("/v1/users/:id", authMiddleware, updateUser(s))
	f.Delete("/v1/users/:id", authMiddleware, deleteUser(s))

	f.Post("/v1/customers", createCustomer(s))
	f.Get("/v1/customers", listCustomers(s))
	f.Delete("v1/customers/:id", deleteCustomer(s))

	f.Get("/v1/products", listProducts(s))
	f.Get("/v1/products/:id", findProduct(s))

	f.Post("/v1/products", authMiddleware, addProduct(s))
	f.Put("/v1/products/:id", authMiddleware, updateProduct(s))
	f.Delete("/v1/products/:id", authMiddleware, deleteProduct(s))

	f.Post("/v1/invoices", authMiddleware, generateInvoice(s))

	f.Get("/swagger/*", swagger.WrapHandler)

	return f
}

// authMiddleware for handlers that require user login.
func authMiddleware(c *fiber.Ctx) error {
	token := c.Request().Header.Peek("Authorization")
	_, err := jwt.Verify(string(token))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "you don't have authorization"})
	}
	return c.Next()
}

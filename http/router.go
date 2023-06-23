package http

import (
	"net/http"

	"github.com/adrianolmedo/genesis/app"
	"github.com/adrianolmedo/genesis/http/jwt"
	"github.com/adrianolmedo/genesis/postgres"

	"github.com/gofiber/fiber/v2"
)

// Router call all routes with its handlers.
func Router(strg *postgres.Storage) *fiber.App {
	s := app.NewServices(strg)
	f := fiber.New()

	g := f.Group("/v1/users")
	g.Use(authMiddleware)
	g.Get("", listUsers(s))
	g.Put("/:id", updateUser(s))
	g.Delete("/:id", deleteUser(s))

	f.Post("v1/login", loginUser(s))
	f.Post("v1/users", signUpUser(s))
	f.Get("/v1/users/:id", findUser(s))

	f.Post("/v1/customers", createCustomer(s))
	f.Get("/v1/customers", listCustomers(s))

	f.Get("/v1/products", listProducts(s))
	f.Get("/v1/products/:id", findProduct(s))

	f.Post("/v1/products", addProduct(s), authMiddleware)
	f.Put("/v1/products/:id", updateProduct(s), authMiddleware)
	f.Delete("/v1/products/:id", deleteProduct(s), authMiddleware)

	f.Post("v1/invoices", generateInvoice(s), authMiddleware)
	return f
}

// authMiddleware for handlers that require user login.
func authMiddleware(c *fiber.Ctx) error {
	token := c.Request().Header.Peek("Authorization")
	_, err := jwt.Verify(string(token))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message_error": "you don't have authorization"})
	}
	return c.Next()
}

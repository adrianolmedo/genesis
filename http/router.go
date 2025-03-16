package http

import (
	"net/http"

	"github.com/adrianolmedo/genesis/app"
	_ "github.com/adrianolmedo/genesis/docs"
	"github.com/adrianolmedo/genesis/http/jwt"
	"github.com/adrianolmedo/genesis/logger"
	"github.com/adrianolmedo/genesis/pgsql/pq"

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
func Router(strg *pq.Storage) *fiber.App {
	s := app.NewServices(strg)
	f := fiber.New(fiber.Config{
		DisableStartupMessage: true, // Disables Fiber's startup message
	})

	f.Get("/v1/test", func(c *fiber.Ctx) error {
		logger.Info("testing", "path", c.Path())
		return successJSON(c, http.StatusOK, respDetails{
			Message: "The server is ok",
		})
	})
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

// successJSON respond JSON.
func successJSON(c *fiber.Ctx, httpStatus int, details respDetails) error {
	return c.Status(http.StatusCreated).JSON(successResp{
		Status:      "success",
		respDetails: details,
	})
}

// errorJSON respond JSON.
func errorJSON(c *fiber.Ctx, httpStatus int, details respDetails) error {
	return c.Status(http.StatusCreated).JSON(errorResp{
		Status: "error",
		Error:  details,
	})
}

// authWare middleware for handlers that require user login.
func authWare(c *fiber.Ctx) error {
	token := c.Request().Header.Peek("Authorization")
	_, err := jwt.Verify(string(token))
	if err != nil {
		return errorJSON(c, http.StatusBadRequest, respDetails{
			Code:    "001",
			Message: "You don't have authorization",
			Details: "Sign to access",
		})
	}
	return c.Next()
}

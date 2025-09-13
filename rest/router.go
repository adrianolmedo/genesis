package rest

import (
	"net/http"
	"time"

	"github.com/adrianolmedo/genesis/compose"
	_ "github.com/adrianolmedo/genesis/docs"
	"github.com/adrianolmedo/genesis/rest/jwt"

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
func Router(svcs *compose.Services) *fiber.App {
	f := fiber.New()
	rateLimit := newRateLimit(2, 5, 5*time.Minute) // 2 req/sec, burst of 5, cleanup inactive IPs after 5 min
	f.Use(rateLimit.middlewarePerIP, timeoutWare(60*time.Second))
	f.Get("/v1/test", func(c *fiber.Ctx) error {
		return respJSON(c, http.StatusOK, detailsResp{
			Message: "Hello world",
		})
	})
	f.Get("/v1/test-ratelimit", testRatelimit())
	f.Get("/v1/test-timeout", testTimeout())
	f.Post("/v1/login", loginUser(svcs))
	f.Post("/v1/users", signUpUser(svcs))
	f.Get("/v1/users/:id", findUser(svcs))
	f.Get("/v1/users", authWare, listUsers(svcs))
	f.Put("/v1/users/:id", authWare, updateUser(svcs))
	f.Delete("/v1/users/:id", authWare, deleteUser(svcs))
	f.Post("/v1/customers", createCustomer(svcs))
	f.Get("/v1/customers", listCustomers(svcs))
	f.Delete("v1/customers/:id", deleteCustomer(svcs))
	f.Get("/v1/products", listProducts(svcs))
	f.Get("/v1/products/:id", findProduct(svcs))
	f.Post("/v1/products", authWare, addProduct(svcs))
	f.Put("/v1/products/:id", authWare, updateProduct(svcs))
	f.Delete("/v1/products/:id", authWare, deleteProduct(svcs))
	f.Post("/v1/invoices", authWare, generateInvoice(svcs))
	f.Get("/swagger/*", swagger.WrapHandler)
	return f
}

// respJSON respond JSON.
func respJSON(c *fiber.Ctx, httpStatus int, details detailsResp) error {
	return c.Status(httpStatus).JSON(resp{
		Status:      "Success",
		detailsResp: details,
	})
}

// errorJSON respond JSON.
func errorJSON(c *fiber.Ctx, httpStatus int, details detailsResp) error {
	return c.Status(httpStatus).JSON(errorResp{
		Status: "Error",
		Error:  details,
	})
}

// authWare middleware for handlers that require user login.
func authWare(c *fiber.Ctx) error {
	token := c.Request().Header.Peek("Authorization")
	_, err := jwt.Verify(string(token))
	if err != nil {
		return errorJSON(c, http.StatusBadRequest, detailsResp{
			Code:    "001",
			Message: "You aren't authenticated",
			Details: "Sign to access",
		})
	}
	return c.Next()
}

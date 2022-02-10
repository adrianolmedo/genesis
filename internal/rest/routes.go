package rest

import (
	"github.com/adrianolmedo/go-restapi/internal/rest/middleware"
	"github.com/adrianolmedo/go-restapi/internal/service"

	"github.com/labstack/echo/v4"
)

// Routes call all routes.
func Routes(e *echo.Echo, s service.Service) {
	e.POST("/v1/login", loginUser(s))
	e.POST("/v1/signup", signUpUser(s))
	e.GET("/v1/users/:id", findUser(s))

	u := e.Group("/v1/users")
	u.Use(middleware.Auth)  // Routes that required authentication.
	u.GET("", listUsers(s)) // E.g.: GET /v1/users
	u.PUT("/:id", updateUser(s))
	u.DELETE("/:id", deleteUser(s))

	e.GET("/v1/products", listProducts(s))
	e.GET("/v1/products/:id", findProduct(s))

	p := e.Group("/v1/products")
	p.Use(middleware.Auth)
	p.POST("", addProduct(s))
	p.PUT("/:id", updateProduct(s))
	p.DELETE("/:id", deleteProduct(s))

	i := e.Group("/v1/invoices")
	i.POST("", generateInvoice(s))
}

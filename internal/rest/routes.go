package rest

import (
	"github.com/adrianolmedo/genesis/internal/rest/middleware"
	"github.com/adrianolmedo/genesis/internal/service"

	"github.com/labstack/echo/v4"
)

// Routes call all routes.
func Routes(e *echo.Echo, s service.Service) {
	e.POST("/v1/login", loginUser(s))
	e.POST("/v1/signup", signUpUser(s))
	e.GET("/v1/users/:id", findUser(s))

	u := e.Group("/v1/users")
	u.Use(middleware.Auth) // Routes that required authentication.
	u.GET("", listUsers(s))
	u.PUT("/:id", updateUser(s))
	u.DELETE("/:id", deleteUser(s))

	e.GET("/v1/products", listProducts(s))
	e.GET("/v1/products/:id", findProduct(s))

	e.POST("/v1/products", addProduct(s), middleware.Auth)
	e.PUT("/v1/products/:id", updateProduct(s), middleware.Auth)
	e.DELETE("/v1/products/:id", deleteProduct(s), middleware.Auth)

	i := e.Group("/v1/invoices")
	i.Use(middleware.Auth)
	i.POST("", generateInvoice(s))
}

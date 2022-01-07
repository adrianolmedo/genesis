package rest

import (
	"github.com/adrianolmedo/go-restapi-practice/internal/rest/middleware"
	"github.com/adrianolmedo/go-restapi-practice/internal/service"

	"github.com/labstack/echo/v4"
)

// Routes call all routes.
func Routes(e *echo.Echo, s service.Service) {
	e.POST("/v1/login", loginUser(s))
	e.POST("/v1/signup", signUpUser(s))
	e.GET("/v1/users/:id", findUser(s))

	u := e.Group("/v1/users")
	// Routes that required authentication.
	u.Use(middleware.Auth) // E.g.: GET /v1/users
	u.GET("", listUsers(s))
	u.PUT("/:id", updateUser(s))
	u.DELETE("/:id", deleteUser(s))
}

package rest

import (
	"github.com/adrianolmedo/go-restapi-practice/internal/rest/middleware"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage"

	"github.com/labstack/echo/v4"
)

// Routes call all routes.
// TO-DO: Not pass nil as argument.
func Routes(e *echo.Echo, r *storage.Repositories) {
	e.POST("/v1/login", loginUser(r.LoginRepository))
	e.POST("/v1/signup", signUpUser(r.UserRepository))
	e.GET("/v1/users/:id", findUser(r.UserRepository))

	u := e.Group("/v1/users")
	// Routes that required authentication.
	u.Use(middleware.Auth) // E.g.: GET /v1/users
	u.GET("", listUsers(r.UserRepository))
	u.PUT("/:id", updateUser(r.UserRepository))
	u.DELETE("/:id", deleteUser(r.UserRepository))
}

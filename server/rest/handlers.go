package rest

import (
	"github.com/adrianolmedo/go-restapi-practice/user"

	"github.com/labstack/echo/v4"
)

// Handlers without Auth middleware.
func Handlers(e *echo.Echo, u user.Repository) {
	e.POST("/v1/login", login)
	e.POST("/v1/signup", createUser(u))
	e.GET("/v1/users/:id", getUserByID(u))
}

// HandlersAuthRequired end-points that required JWT authentication.
func HandlersAuthRequired(e *echo.Echo, r user.Repository) {
	u := e.Group("/v1/users")
	u.Use(Auth)
	u.GET("", getAllUsers(r)) // GET /v1/users
	u.PUT("/:id", updateUserByID(r))
	u.DELETE("/:id", deleteUserByID(r))
}

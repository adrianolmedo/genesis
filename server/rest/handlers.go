package rest

import (
	"go-restapi-practice/service"

	"github.com/labstack/echo/v4"
)

// Handlers without Auth middleware.
func Handlers(e *echo.Echo, u service.UserDAO) {
	e.POST("/v1/login", login)
	e.POST("/v1/signup", createUser(u))
	e.GET("/v1/users/:id", getUserByID(u))
}

// HandlersAuthRequired end-points that required JWT authentication.
func HandlersAuthRequired(e *echo.Echo, u service.UserDAO) {
	user := e.Group("/v1/users")
	user.Use(Auth)
	user.GET("", getAllUsers(u)) // GET /v1/users
	user.PUT("/:id", updateUserByID(u))
	user.DELETE("/:id", deleteUserByID(u))
}

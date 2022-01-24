package middleware

import (
	"net/http"

	"github.com/adrianolmedo/go-restapi-practice/jwt"

	"github.com/labstack/echo/v4"
)

// Auth middleware for check JWT authentication.
func Auth(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		_, err := jwt.Verify(token)
		if err != nil {
			return c.JSON(http.StatusForbidden, map[string]string{"message_error": "you don't have authorization"})
		}
		return f(c)
	}
}

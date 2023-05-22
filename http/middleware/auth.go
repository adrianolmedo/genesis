package middleware

import (
	"net/http"

	"github.com/adrianolmedo/go-restapi/http/jwt"

	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	token := c.Request().Header.Peek("Authorization")
	_, err := jwt.Verify(string(token))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message_error": "you don't have authorization"})
	}
	return c.Next()
}

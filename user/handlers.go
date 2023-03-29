package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/adrianolmedo/go-restapi/api"

	"github.com/gofiber/fiber/v2"
)

// findUser handler GET: /users/:id
func findUser(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := api.RespJSON(api.MsgError, "Positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		user, err := s.Find(int64(id))
		if errors.Is(err, ErrUserNotFound) {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusNotFound).JSON(resp)
		}

		if err != nil {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		resp := api.RespJSON(api.MsgOK, "", UserProfileDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
		return c.Status(http.StatusOK).JSON(resp)
	}
}

package store

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// addProduct handler POST: /products
func addProduct(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := AddProductForm{}

		err := c.BodyParser(&form)
		if err != nil {
			resp := newResponse(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.Add(&Product{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})

		if err != nil {
			resp := newResponse(msgError, "", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		resp := newResponse(msgOK, "product added", ProductCardDTO{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})

		return c.Status(http.StatusCreated).JSON(resp)
	}
}

package product

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// addProduct handler POST: /products
func addProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := AddProductForm{}

		err := c.BodyParser(&form)
		if err != nil {
			resp := newResponse(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		resp := newResponse(msgOK, "product added", ProductCardDTO{
			Name:         "Creatina",
			Observations: "Muscle growth",
			Price:        100,
		})

		return c.Status(http.StatusCreated).JSON(resp)
	}
}

package store

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrianolmedo/go-restapi/api"

	"github.com/gofiber/fiber/v2"
)

// addProduct handler POST: /products
func addProduct(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := AddProductForm{}
		err := c.BodyParser(&form)
		if err != nil {
			resp := api.RespJSON(api.MsgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.Add(&Product{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})

		if err != nil {
			resp := api.RespJSON(api.MsgError, "", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Add logger message: "New product added"

		resp := api.RespJSON(api.MsgOK, "product added", ProductCardDTO{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})

		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// listProducts handler GET: /products
func listProducts(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		products, err := s.List()
		if err != nil {
			resp := api.RespJSON(api.MsgError, "", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		if len(products) == 0 {
			resp := api.RespJSON(api.MsgOK, "there are not products", nil)
			return c.Status(http.StatusOK).JSON(resp)
		}

		list := make(ProductList, 0, len(products))

		assemble := func(p *Product) ProductCardDTO {
			return ProductCardDTO{
				ID:           p.ID,
				Name:         p.Name,
				Observations: p.Observations,
				Price:        p.Price,
			}
		}

		for _, v := range products {
			list = append(list, assemble(v))
		}

		resp := api.RespJSON(api.MsgOK, "", list)
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// findProduct handler GET: /products/:id
func findProduct(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := api.RespJSON(api.MsgError, "positive number expected for ID product", nil)
			return c.Status(http.StatusBadRequest).JSON(resp) // 400
		}

		product, err := s.Find(int64(id))
		if errors.Is(err, ErrProductNotFound) {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusNotFound).JSON(resp) // 404
		}

		if err != nil {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		resp := api.RespJSON(api.MsgOK, "", ProductCardDTO{
			ID:           product.ID,
			Name:         product.Name,
			Observations: product.Observations,
			Price:        product.Price,
		})
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// updateProduct handler PUT: /products/:id
func updateProduct(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		// TO-DO: Add logger message: "Request to update product ID %d"

		if id < 0 || err != nil {
			resp := api.RespJSON(api.MsgError, "positive number expected for ID product", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		form := UpdateProductForm{}
		err = c.BodyParser(&form)
		if err != nil {
			resp := api.RespJSON(api.MsgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		form.ID = int64(id)

		err = s.Update(Product{
			ID:           form.ID,
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})
		if errors.Is(err, ErrProductNotFound) {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Add logger message: "Product ID %d updated"

		resp := api.RespJSON(api.MsgOK, "product updated", form)
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// deleteProduct handler DELETE: /products/:id
func deleteProduct(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := api.RespJSON(api.MsgError, "positive number expected for ID product", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.Remove(int64(id))
		if errors.Is(err, ErrProductNotFound) {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := api.RespJSON(api.MsgError, fmt.Sprintf("could not delete product: %s", err), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Add logger mesaage: "Product with ID %d removed from DB"

		resp := api.RespJSON(api.MsgOK, "product deleted", nil)
		return c.Status(http.StatusOK).JSON(resp) // maybe 204
	}
}

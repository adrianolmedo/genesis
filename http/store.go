package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/app"

	"github.com/gofiber/fiber/v2"
)

// addProduct handler POST: /products
func addProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := domain.AddProductForm{}
		err := c.BodyParser(&form)
		if err != nil {
			resp := respJSON(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.Store.Add(&domain.Product{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})

		if err != nil {
			resp := respJSON(msgError, "", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Add logger message: "New product added"

		resp := respJSON(msgOK, "product added", domain.ProductCardDTO{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})

		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// listProducts handler GET: /products
func listProducts(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		products, err := s.Store.List()
		if err != nil {
			resp := respJSON(msgError, "", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		if products.IsEmpty() {
			resp := respJSON(msgOK, "there are not products", nil)
			return c.Status(http.StatusOK).JSON(resp)
		}

		list := make([]domain.ProductCardDTO, 0, len(products))

		assemble := func(p *domain.Product) domain.ProductCardDTO {
			return domain.ProductCardDTO{
				ID:           p.ID,
				Name:         p.Name,
				Observations: p.Observations,
				Price:        p.Price,
			}
		}

		for _, v := range products {
			list = append(list, assemble(v))
		}

		resp := respJSON(msgOK, "", list)
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// createCustomer handler POST: /customer
func createCustomer(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := domain.CreateCustomerForm{}
		err := c.BodyParser(&form)
		if err != nil {
			resp := respJSON(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.Store.AddCustomer(&domain.Customer{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
		})

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		resp := respJSON(msgOK, "customer created", domain.CustomerProfileDTO{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
		})

		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// listCustomers handler GET: /customers
func listCustomers(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		f, err := getFilter(c)
		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp) // 400
		}

		fr, err := s.Store.ListCustomers(f)
		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		var firstPage, lastPage, previousPage, nextPage string

		// Set first page and last page for the pagination reponse.
		firstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", c.Path(), f.Limit, 0, f.Sort)
		lastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", c.Path(), f.Limit, fr.TotalPages, f.Sort)

		// Set previous page pagination response.
		if f.Page > 0 {
			previousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", c.Path(), f.Limit, f.Page-1, f.Sort)
		}

		// Set next pagination response.
		if f.Page < fr.TotalPages {
			nextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", c.Path(), f.Limit, f.Page+1, f.Sort)
		}

		// Reset previous page.
		if f.Page > fr.TotalPages {
			previousPage = ""
		}

		links := domain.PaginationLinks{
			Page:         f.Page,
			Limit:        f.Limit,
			Sort:         f.Sort,
			TotalRows:    fr.TotalRows,
			TotalPages:   fr.TotalPages,
			FirstPage:    firstPage,
			PreviousPage: previousPage,
			NextPage:     nextPage,
			LastPage:     lastPage,
			FromRow:      fr.FromRow,
			ToRow:        fr.ToRow,
		}

		customers, ok := fr.Rows.(domain.Customers)
		if !ok {
			resp := respJSON(msgError, "error data assertion", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		if customers.IsEmpty() {
			resp := respJSON(msgOK, "there are not customers", nil)
			return c.Status(http.StatusOK).JSON(resp)
		}

		data := make([]domain.CustomerProfileDTO, 0, len(customers))

		assemble := func(cx *domain.Customer) domain.CustomerProfileDTO {
			return domain.CustomerProfileDTO{
				FirstName: cx.FirstName,
				LastName:  cx.LastName,
				Email:     cx.Email,
			}
		}

		for _, v := range customers {
			data = append(data, assemble(v))
		}

		resp := respJSON(msgOK, "", data).links(links)
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// findProduct handler GET: /products/:id
func findProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for ID product", nil)
			return c.Status(http.StatusBadRequest).JSON(resp) // 400
		}

		product, err := s.Store.Find(int64(id))
		if errors.Is(err, domain.ErrProductNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNotFound).JSON(resp) // 404
		}

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		resp := respJSON(msgOK, "", domain.ProductCardDTO{
			ID:           product.ID,
			Name:         product.Name,
			Observations: product.Observations,
			Price:        product.Price,
		})
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// updateProduct handler PUT: /products/:id
func updateProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		// TO-DO: Add logger message: "Request to update product ID %d"

		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for ID product", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		form := domain.UpdateProductForm{}
		err = c.BodyParser(&form)
		if err != nil {
			resp := respJSON(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		form.ID = int64(id)

		err = s.Store.Update(domain.Product{
			ID:           form.ID,
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})
		if errors.Is(err, domain.ErrProductNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Add logger message: "Product ID %d updated"

		resp := respJSON(msgOK, "product updated", form)
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// deleteProduct handler DELETE: /products/:id
func deleteProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for ID product", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.Store.Remove(int64(id))
		if errors.Is(err, domain.ErrProductNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, fmt.Sprintf("could not delete product: %s", err), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Add logger mesaage: "Product with ID %d removed from DB"

		resp := respJSON(msgOK, "product deleted", nil)
		return c.Status(http.StatusOK).JSON(resp) // maybe 204
	}
}

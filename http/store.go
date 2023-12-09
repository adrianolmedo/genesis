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

// addProduct godoc
// @Summary	Add product
// @Description	Register product
// @Tags products
// @Accept json
// @Produce	json
// @Failure	400	{object}	respError
// @Failure	401	{object}	respError
// @Failure	500	{object}	respError
// @Success	201	{object}	respOkData{data=productCardDTO}
// @Param addProductForm	body addProductForm true "application/json"
// @Router /products [post]
func addProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := addProductForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{"the JSON structure is not correct"})
		}

		err = s.Store.Add(&domain.Product{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(respError{err.Error()})
		}

		// TODO: Add logger message: "New product added"

		return c.Status(http.StatusCreated).JSON(respOkData{
			Msg: "product added",
			Data: productCardDTO{
				Name:         form.Name,
				Observations: form.Observations,
				Price:        form.Price,
			},
		})
	}
}

// addProductForm represents a subset of fields to create a Product.
type addProductForm struct {
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

// productCardDTO subset of Product fields.
type productCardDTO struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
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

		list := make([]productCardDTO, 0, len(products))

		assemble := func(p *domain.Product) productCardDTO {
			return productCardDTO{
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
			Password:  form.Password,
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

// deleteCustomer handler DELETE: /customers/:id
func deleteCustomer(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for customer ID", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.Store.RemoveCustomer(id)
		if errors.Is(err, domain.ErrProductNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, fmt.Sprintf("could not delete customer: %s", err), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Add logger mesaage: "Customer with ID %d removed from DB"

		resp := respJSON(msgOK, "customer deleted", nil)
		return c.Status(http.StatusOK).JSON(resp) // maybe 204
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

		customers, ok := fr.Rows.(domain.Customers)
		if !ok {
			resp := respJSON(msgError, "error in data assertion", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		if customers.IsEmpty() {
			resp := respJSON(msgOK, "there are not customers", nil)
			return c.Status(http.StatusOK).JSON(resp)
		}

		data := make([]domain.CustomerProfileDTO, 0, len(customers))

		assemble := func(cx *domain.Customer) domain.CustomerProfileDTO {
			return domain.CustomerProfileDTO{
				ID:        cx.ID,
				FirstName: cx.FirstName,
				LastName:  cx.LastName,
				Email:     cx.Email,
			}
		}

		for _, v := range customers {
			data = append(data, assemble(v))
		}

		ls := f.GenLinksResp(c.Path(), fr.TotalPages)
		resp := respJSON(msgOK, "", data).setLinks(ls).setMeta(fr)
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

		product, err := s.Store.Find(id)
		if errors.Is(err, domain.ErrProductNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNotFound).JSON(resp) // 404
		}

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		resp := respJSON(msgOK, "", productCardDTO{
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

		form.ID = id

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

		err = s.Store.Remove(id)
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

package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/app"
	"github.com/adrianolmedo/genesis/pgsql"

	"github.com/gofiber/fiber/v2"
)

// addProduct godoc
//
//	@Summary		Add product
//	@Description	Register product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400				{object}	respError
//	@Failure		401				{object}	respError
//	@Failure		500				{object}	respError
//	@Success		201				{object}	respOkData{data=productCardDTO}
//	@Param			addProductForm	body		addProductForm	true	"application/json"
//	@Router			/products [post]
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

// listProduct godoc
//
//	@Summary		List products
//	@Description	Get a list of products
//	@Tags			products
//	@Produce		json
//	@Failure		500	{object}	respError
//	@Success		200	{object}	respOk
//	@Success		200	{object}	respData{data=[]productCardDTO}
//	@Router			/products [get]
func listProducts(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		products, err := s.Store.List()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(respError{err.Error()})
		}

		if products.IsEmpty() {
			return c.Status(http.StatusOK).JSON(respOk{"there are not products"})
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

		return c.Status(http.StatusOK).JSON(respData{
			Data: list,
		})
	}
}

// createCustomer godoc
//
//	@Summary		Create customer
//	@Description	Set new customer
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400					{object}	respError
//	@Failure		500					{object}	respError
//	@Success		201					{object}	respOkData{data=customerProfileDTO}
//	@Param			createCustomerForm	body		createCustomerForm	true	"application/json"
//	@Router			/customer [post]
func createCustomer(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := createCustomerForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{
				"the JSON structure is not correct",
			})
		}

		err = s.Store.AddCustomer(&domain.Customer{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(respError{err.Error()})
		}

		return c.Status(http.StatusCreated).JSON(respOkData{
			Msg: "customer created",
			Data: customerProfileDTO{
				FirstName: form.FirstName,
				LastName:  form.LastName,
				Email:     form.Email,
			},
		})
	}
}

// customerProfileDTO subset of Customer fields.
type customerProfileDTO struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// createCustomerForm subset of fields to request to create a Customer.
type createCustomerForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// deleteCustomer godoc
//
//	@Summary		Create customer
//	@Description	Set new customer
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	respError
//	@Failure		204	{object}	respError
//	@Failure		500	{object}	respError
//	@Success		200	{object}	respOkData{data=customerProfileDTO}
//	@Param			id	path		int	true	"Customer id"
//	@Router			/customers/{id} [delete]
func deleteCustomer(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{
				"positive number expected for customer ID",
			})
		}

		err = s.Store.RemoveCustomer(id)
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(http.StatusNoContent).JSON(respError{err.Error()})
		}

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(respError{
				fmt.Sprintf("could not delete customer: %s", err),
			})
		}

		// TO-DO: Add logger message: "Customer with ID %d removed from DB"
		return c.Status(http.StatusOK).JSON(respOk{"customer deleted"}) // maybe 204
	}
}

// listCustomers godoc
//
//	@Summary		List customers
//	@Description	Paginate customers
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400			{object}	respError
//	@Failure		500			{object}	respError
//	@Success		200			{object}	respMetaData{links=genesis.LinksResp,meta=genesis.FilteredResults,data=[]customerProfileDTO}
//	@Param			limit		query		int		false	"Limit of pages"					example(2)
//	@Param			page		query		int		false	"Current page"						example(1)
//	@Param			sort		query		string	false	"Sort results by a value"			example(created_at)
//	@Param			direction	query		string	false	"Order by ascendent o descendent"	example(desc)
//	@Router			/customers [get]
func listCustomers(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		p, err := pgsql.NewPager(
			c.QueryInt("limit"),
			c.QueryInt("page"),
			c.Query("sort", "created_at"),
			c.Query("direction"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{err.Error()}) // 400
		}

		pr, err := s.Store.ListCustomers(p)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(respError{err.Error()})
		}

		customers, ok := pr.Rows.(domain.Customers)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(respError{"data assertion"})
		}

		if customers.IsEmpty() {
			return c.Status(http.StatusOK).JSON(respOk{"there are not customers"})
		}

		assemble := func(cx *domain.Customer) customerProfileDTO {
			return customerProfileDTO{
				ID:        cx.ID,
				FirstName: cx.FirstName,
				LastName:  cx.LastName,
				Email:     cx.Email,
			}
		}

		data := make([]customerProfileDTO, 0, len(customers))
		for _, v := range customers {
			data = append(data, assemble(v))
		}

		return c.Status(http.StatusOK).JSON(respMetaData{
			Links: p.GenLinks(c.Path(), pr.TotalPages),
			Meta:  pr,
			Data:  data,
		})
	}
}

// findProduct godoc
//
//	@Summary		Find product
//	@Description	Find product by its id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	respError
//	@Failure		404	{object}	respError
//	@Success		200	{object}	respData{data=productCardDTO}
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [get]
func findProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{
				"positive number expected for ID product",
			}) // 400
		}

		product, err := s.Store.Find(id)
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(http.StatusNotFound).JSON(respError{err.Error()}) // 404
		}

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{err.Error()})
		}

		return c.Status(http.StatusOK).JSON(respData{
			Data: productCardDTO{
				ID:           product.ID,
				Name:         product.Name,
				Observations: product.Observations,
				Price:        product.Price,
			},
		})
	}
}

// updateProduct godoc
//
//	@Summary		Update product
//	@Description	Update product by its id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	respError
//	@Failure		204	{object}	respError
//	@Failure		500	{object}	respError
//	@Success		200	{object}	respOk
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [put]
func updateProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		// TO-DO: Add logger message: "Request to update product ID %d"

		if id < 0 || err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{
				"positive number expected for ID product",
			})
		}

		form := updateProductForm{}
		err = c.BodyParser(&form)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{
				"the JSON structure is not correct",
			})
		}

		form.ID = id

		err = s.Store.Update(domain.Product{
			ID:           form.ID,
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(http.StatusNoContent).JSON(respError{err.Error()})
		}

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(respError{err.Error()})
		}

		// TO-DO: Add logger message: "Product ID %d updated"
		return c.Status(http.StatusOK).JSON(respOk{"product updated"})
	}
}

// updateProductForm represents a subset of fields to update a Product.
type updateProductForm struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Observations string  `json:"observations"`
	Price        float64 `json:"price"`
}

// deleteProduct godoc
//
//	@Summary		Delete product
//	@Description	Delete product by its id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	respError
//	@Failure		204	{object}	respError
//	@Failure		500	{object}	respError
//	@Success		200	{object}	respOk
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [delete]
func deleteProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{
				"positive number expected for ID product",
			})
		}

		err = s.Store.Remove(id)
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(http.StatusNoContent).JSON(respError{err.Error()})
		}

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(respError{
				fmt.Sprintf("could not delete product: %s", err),
			})
		}

		// TO-DO: Add logger mesaage: "Product with ID %d removed from DB"
		return c.Status(http.StatusOK).JSON(respOk{"product deleted"}) // maybe 204
	}
}

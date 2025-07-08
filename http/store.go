package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/app"
	"github.com/adrianolmedo/genesis/logger"
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
//	@Failure		400				{object}	errorResp
//	@Failure		401				{object}	errorResp
//	@Failure		500				{object}	errorResp
//	@Success		201				{object}	successResp{data=productCardDTO}
//	@Param			addProductForm	body		addProductForm	true	"application/json"
//	@Router			/products [post]
func addProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := addProductForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		product := &domain.Product{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		}

		err = s.Store.Add(product)

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		logger.Debug("Product added", "product", product.UUID)

		return successJSON(c, http.StatusCreated, respDetails{
			Message: "Product added",
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
	Name         string `json:"name"`
	Observations string `json:"observations"`
	Price        int64  `json:"price"`
}

// productCardDTO subset of Product fields.
type productCardDTO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Observations string `json:"observations"`
	Price        int64  `json:"price"`
}

// listProduct godoc
//
//	@Summary		List products
//	@Description	Get a list of products
//	@Tags			products
//	@Produce		json
//	@Failure		500	{object}	errorResp
//	@Success		200	{object}	successResp
//	@Success		200	{object}	successResp{data=[]productCardDTO}
//	@Router			/products [get]
func listProducts(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		products, err := s.Store.List()
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if products.IsEmpty() {
			return successJSON(c, http.StatusOK, respDetails{
				Code:    "005",
				Message: "There are not products",
			})
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

		return successJSON(c, http.StatusOK, respDetails{
			Message: "Ok",
			Data:    list,
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
//	@Failure		400					{object}	errorResp
//	@Failure		500					{object}	errorResp
//	@Success		201					{object}	successResp{data=customerProfileDTO}
//	@Param			createCustomerForm	body		createCustomerForm	true	"application/json"
//	@Router			/customer [post]
func createCustomer(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := createCustomerForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		err = s.Store.AddCustomer(&domain.Customer{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		return successJSON(c, http.StatusCreated, respDetails{
			Message: "Customer created",
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
	ID        int64  `json:"id,omitempty"`
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
//	@Failure		400	{object}	errorResp
//	@Failure		204	{object}	errorResp
//	@Failure		500	{object}	errorResp
//	@Success		200	{object}	successResp{data=customerProfileDTO}
//	@Param			id	path		int	true	"Customer id"
//	@Router			/customers/{id} [delete]
func deleteCustomer(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}

		err = s.Store.RemoveCustomer(int64(id))
		if errors.Is(err, domain.ErrProductNotFound) {
			return errorJSON(c, http.StatusNoContent, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: "Coudl not delete customer",
			})
		}

		logger.Info("customer", fmt.Sprintf("customer with ID %d removed from DB", id))

		return successJSON(c, http.StatusOK, respDetails{
			Message: "Customer delete",
		})
	}
}

// listCustomers godoc
//
//	@Summary		List customers
//	@Description	Paginate customers
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400			{object}	errorResp
//	@Failure		500			{object}	errorResp
//	@Success		200			{object}	pagerResp{links=pgsql.PagerLinks,meta=pgsql.PagerResults,data=[]customerProfileDTO}
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
			c.Query("direction"),
		)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		pr, err := s.Store.ListCustomers(p)
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		customers, ok := pr.Rows.(domain.Customers)
		if !ok {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: "Data assertion",
			})
		}

		if customers.IsEmpty() {
			return successJSON(c, http.StatusOK, respDetails{
				Message: "There are not customers",
			})
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

		return c.Status(http.StatusOK).JSON(pagerResp{
			Links: p.Links(c.Path(), pr.TotalPages),
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
//	@Failure		400	{object}	errorResp
//	@Failure		404	{object}	errorResp
//	@Success		200	{object}	successResp{data=productCardDTO}
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [get]
func findProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID product",
			})
		}

		product, err := s.Store.Find(int64(id))
		if errors.Is(err, domain.ErrProductNotFound) {
			return errorJSON(c, http.StatusNotFound, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		return successJSON(c, http.StatusOK, respDetails{
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
//	@Failure		400	{object}	errorResp
//	@Failure		204	{object}	errorResp
//	@Failure		500	{object}	errorResp
//	@Success		200	{object}	successResp
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [put]
func updateProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		logger.Debug("product", fmt.Sprintf("request to update product ID %d", id))

		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID product",
			})
		}

		form := updateProductForm{}
		err = c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		form.ID = int64(id)

		err = s.Store.Update(domain.Product{
			ID:           form.ID,
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})
		if errors.Is(err, domain.ErrProductNotFound) {
			return errorJSON(c, http.StatusNoContent, respDetails{
				Code:    "002",
				Message: err.Error(),
			})
		}

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "002",
				Message: err.Error(),
			})
		}

		logger.Debug("product", fmt.Sprintf("product ID %d updated", id))

		return successJSON(c, http.StatusOK, respDetails{
			Message: "Product updated",
		})

	}
}

// updateProductForm represents a subset of fields to update a Product.
type updateProductForm struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Observations string `json:"observations"`
	Price        int64  `json:"price"`
}

// deleteProduct godoc
//
//	@Summary		Delete product
//	@Description	Delete product by its id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	errorResp
//	@Failure		204	{object}	errorResp
//	@Failure		500	{object}	errorResp
//	@Success		200	{object}	successResp
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [delete]
func deleteProduct(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID product",
			})
		}

		err = s.Store.Remove(int64(id))
		if errors.Is(err, domain.ErrProductNotFound) {
			return errorJSON(c, http.StatusNoContent, respDetails{
				Message: err.Error(),
			})
		}

		if err != nil {
			return errorJSON(c, http.StatusNoContent, respDetails{
				Message: fmt.Sprintf("Could not delete product: %s", err),
			})
		}

		logger.Debug("product", fmt.Sprintf("product with ID %d removed from DB", id))

		return successJSON(c, http.StatusOK, respDetails{
			Message: "Product deleted",
		})
	}
}

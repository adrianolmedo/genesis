package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrianolmedo/genesis/compose"
	"github.com/adrianolmedo/genesis/logger"
	"github.com/adrianolmedo/genesis/pgsql"
	"github.com/adrianolmedo/genesis/store"

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
//	@Success		201				{object}	resp{data=productCardResp}
//	@Param			addProductReq	body		addProductReq	true	"application/json"
//	@Router			/products [post]
func addProduct(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := addProductReq{}

		err := c.BodyParser(&req)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		product := &store.Product{
			Name:         req.Name,
			Observations: req.Observations,
			Price:        req.Price,
		}

		err = svcs.Store.Add(ctx, product)

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		logger.Debug("Product added", "product", product.UUID)

		return respJSON(c, http.StatusCreated, respDetails{
			Message: "Product added",
			Data: productCardResp{
				Name:         req.Name,
				Observations: req.Observations,
				Price:        req.Price,
			},
		})
	}
}

// addProductReq represents a subset of fields to create a Product.
type addProductReq struct {
	Name         string `json:"name"`
	Observations string `json:"observations"`
	Price        int64  `json:"price"`
}

// productCardResp subset of Product fields.
type productCardResp struct {
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
//	@Success		200	{object}	resp
//	@Success		200	{object}	resp{data=[]productCardResp}
//	@Router			/products [get]
func listProducts(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		products, err := svcs.Store.List(ctx)
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if products.IsEmpty() {
			return respJSON(c, http.StatusOK, respDetails{
				Code:    "005",
				Message: "There are not products",
			})
		}

		list := make([]productCardResp, 0, len(products))
		assemble := func(p store.Product) productCardResp {
			return productCardResp{
				ID:           p.ID,
				Name:         p.Name,
				Observations: p.Observations,
				Price:        p.Price,
			}
		}

		for _, v := range products {
			list = append(list, assemble(v))
		}

		return respJSON(c, http.StatusOK, respDetails{
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
//	@Success		201					{object}	resp{data=customerProfileResp}
//	@Param			createCustomerReq	body		createCustomerReq	true	"application/json"
//	@Router			/customer [post]
func createCustomer(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := createCustomerReq{}

		err := c.BodyParser(&req)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		err = svcs.Store.AddCustomer(ctx, &store.Customer{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		})

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		return respJSON(c, http.StatusCreated, respDetails{
			Message: "Customer created",
			Data: customerProfileResp{
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Email:     req.Email,
			},
		})
	}
}

// customerProfileResp subset of Customer fields.
type customerProfileResp struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// createCustomerReq subset of fields to request to create a Customer.
type createCustomerReq struct {
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
//	@Success		200	{object}	resp{data=customerProfileResp}
//	@Param			id	path		int	true	"Customer id"
//	@Router			/customers/{id} [delete]
func deleteCustomer(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}

		err = svcs.Store.RemoveCustomer(ctx, int64(id))
		if errors.Is(err, store.ErrProductNotFound) {
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

		return respJSON(c, http.StatusOK, respDetails{
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
//	@Success		200			{object}	pagerResp{links=pgsql.PagerLinks,meta=pgsql.PagerResult,data=[]customerProfileResp}
//	@Param			limit		query		int		false	"Limit of pages"					example(2)
//	@Param			page		query		int		false	"Current page"						example(1)
//	@Param			sort		query		string	false	"Sort results by a value"			example(created_at)
//	@Param			direction	query		string	false	"Order by ascendent o descendent"	example(desc)
//	@Router			/customers [get]
func listCustomers(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

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

		pr, err := svcs.Store.ListCustomers(ctx, p)
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		customers, ok := pr.Rows.(store.Customers)
		if !ok {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: "Data assertion",
			})
		}

		if customers.IsEmpty() {
			return respJSON(c, http.StatusOK, respDetails{
				Message: "There are not customers",
			})
		}

		assemble := func(cx store.Customer) customerProfileResp {
			return customerProfileResp{
				ID:        cx.ID,
				FirstName: cx.FirstName,
				LastName:  cx.LastName,
				Email:     cx.Email,
			}
		}

		data := make([]customerProfileResp, 0, len(customers))
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
//	@Success		200	{object}	resp{data=productCardResp}
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [get]
func findProduct(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID product",
			})
		}

		product, err := svcs.Store.Find(ctx, int64(id))
		if errors.Is(err, store.ErrProductNotFound) {
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

		return respJSON(c, http.StatusOK, respDetails{
			Data: productCardResp{
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
//	@Success		200	{object}	resp
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [put]
func updateProduct(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))

		logger.Debug("product", fmt.Sprintf("request to update product ID %d", id))

		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID product",
			})
		}

		req := updateProductReq{}
		err = c.BodyParser(&req)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		req.ID = int64(id)

		err = svcs.Store.Update(ctx, store.Product{
			ID:           req.ID,
			Name:         req.Name,
			Observations: req.Observations,
			Price:        req.Price,
		})
		if errors.Is(err, store.ErrProductNotFound) {
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

		return respJSON(c, http.StatusOK, respDetails{
			Message: "Product updated",
		})

	}
}

// updateProductReq represents a subset of fields to update a Product.
type updateProductReq struct {
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
//	@Success		200	{object}	resp
//	@Param			id	path		int	true	"Product id"
//	@Router			/products/{id} [delete]
func deleteProduct(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID product",
			})
		}

		err = svcs.Store.Remove(ctx, int64(id))
		if errors.Is(err, store.ErrProductNotFound) {
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

		return respJSON(c, http.StatusOK, respDetails{
			Message: "Product deleted",
		})
	}
}

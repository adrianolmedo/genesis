package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/service"

	"github.com/labstack/echo/v4"
)

// addProduct handler POST: /products
func addProduct(s service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		form := domain.AddProductForm{}

		err := c.Bind(&form)
		if err != nil {
			resp := newResponse(msgError, "ER002", "the JSON structure is not correct", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = s.Store.Add(&domain.Product{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})
		if err != nil {
			resp := newResponse(msgError, "ER004", err.Error(), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(msgOK, "OK002", "product added", domain.ProductCardDTO{
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})

		return c.JSON(http.StatusCreated, resp)
	}
}

// findProduct handler GET: /products/:id
func findProduct(s service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if id < 0 || err != nil {
			resp := newResponse(msgError, "ER005", "positive number expected for ID product", nil)
			return c.JSON(http.StatusBadRequest, resp) // 400
		}

		product, err := s.Store.Find(int64(id))
		if errors.Is(err, domain.ErrProductNotFound) {
			resp := newResponse(msgError, "ER007", err.Error(), nil)
			return c.JSON(http.StatusNotFound, resp) // 404
		}

		if err != nil {
			resp := newResponse(msgError, "ER009", err.Error(), nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		resp := newResponse(msgOK, "OK002", "", domain.ProductCardDTO{
			ID:           product.ID,
			Name:         product.Name,
			Observations: product.Observations,
			Price:        product.Price,
		})
		return c.JSON(http.StatusOK, resp)
	}
}

// updateProduct handler PUT: /products/:id
func updateProduct(s service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if id < 0 || err != nil {
			resp := newResponse(msgError, "ER005", "positive number expected for ID product", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		form := domain.UpdateProductForm{}
		err = c.Bind(&form)
		if err != nil {
			resp := newResponse(msgError, "ER002", "the JSON structure is not correct", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		form.ID = int64(id)

		err = s.Store.Update(domain.Product{
			ID:           form.ID,
			Name:         form.Name,
			Observations: form.Observations,
			Price:        form.Price,
		})
		if errors.Is(err, domain.ErrProductNotFound) {
			resp := newResponse(msgError, "ER007", err.Error(), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(msgError, "ER004", err.Error(), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(msgOK, "OK002", "user updated", form)
		return c.JSON(http.StatusOK, resp)
	}
}

// listProducts handler GET: /products.
func listProducts(s service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := s.Store.List()
		if err != nil {
			resp := newResponse(msgError, "ER003", err.Error(), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		if len(products) == 0 {
			resp := newResponse(msgOK, "OK002", "there are not products", nil)
			return c.JSON(http.StatusOK, resp) // maybe 204
		}

		list := make(domain.ProductList, 0, len(products))

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

		resp := newResponse(msgOK, "OK002", "", list)
		return c.JSON(http.StatusOK, resp)
	}
}

// deleteProduct handler DELETE: /products/:id.
func deleteProduct(s service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if id < 0 || err != nil {
			resp := newResponse(msgError, "ER005", "positive number expected for ID product", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = s.Store.Remove(int64(id))
		if errors.Is(err, domain.ErrProductNotFound) {
			resp := newResponse(msgError, "ER007", err.Error(), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(msgError, "ER006", fmt.Sprintf("could not delete product: %s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(msgOK, "OK002", "product deleted", nil)
		return c.JSON(http.StatusOK, resp) // maybe 204
	}
}

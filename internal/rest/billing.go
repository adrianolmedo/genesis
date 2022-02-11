package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/service"

	"github.com/labstack/echo/v4"
)

// generateInvoice handler POST: /invoices
func generateInvoice(s service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		form := domain.GenerateInvoiceForm{}

		err := c.Bind(&form)
		if err != nil {
			resp := newResponse(msgError, "ER002", "the JSON structure is not correct", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		clientID := form.Header.ClientID

		_, err = s.User.Find(clientID)
		if errors.Is(err, domain.ErrUserNotFound) {
			c.Logger().Error("User not found to generate invoice")

			resp := newResponse(msgError, "ER004", err.Error(), nil)
			return c.JSON(http.StatusNotFound, resp) // 404
		}

		if err != nil {
			c.Logger().Error(err)

			resp := newResponse(msgError, "ER009", err.Error(), nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		assemble := func(i *domain.InvoiceItemForm) *domain.InvoiceItem {
			return &domain.InvoiceItem{
				ProductID: i.ProductID,
			}
		}

		items := make(domain.ItemList, 0, len(form.Items))
		for _, item := range form.Items {

			_, err := s.Store.Find(item.ProductID)

			if errors.Is(err, domain.ErrProductNotFound) {
				c.Logger().Error("Product not found to add the invoice")

				resp := newResponse(msgError, "ER007", fmt.Sprintf("%s with id %d", domain.ErrProductNotFound, item.ProductID), nil)
				return c.JSON(http.StatusNotFound, resp) // 404
			}

			if err != nil {
				c.Logger().Error(err)

				resp := newResponse(msgError, "ER009", err.Error(), nil)
				return c.JSON(http.StatusBadRequest, resp)
			}

			items = append(items, assemble(item))

		}

		invoice := &domain.Invoice{
			Header: &domain.InvoiceHeader{
				ClientID: clientID,
			},
			Items: items,
		}

		err = s.Billing.Generate(invoice)
		if err != nil {
			c.Logger().Error(err)

			resp := newResponse(msgError, "ER004", err.Error(), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		c.Logger().Infof("Invoice ID %d generated", invoice.Header.ID)

		resp := newResponse(msgOK, "OK002", "invoice generated", domain.GenerateInvoiceForm{
			Header: form.Header,
			Items:  form.Items,
		})
		return c.JSON(http.StatusCreated, resp)
	}
}

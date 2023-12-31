package http

import (
	"errors"
	"fmt"
	"net/http"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/app"

	"github.com/gofiber/fiber/v2"
)

// generateInvoice godoc
//
//	@Summary		Generate invoice
//	@Description	Generate invoice of products
//	@Tags			billing
//	@Accept			json
//	@Produce		json
//	@Failure		400					{object}	respError
//	@Failure		404					{object}	respError
//	@Failure		500					{object}	respError
//	@Success		201					{object}	respOkData{data=generateInvoiceForm}
//	@Param			generateInvoiceForm	body		generateInvoiceForm	true	"application/json"
//	@Router			/invoices [post]
func generateInvoice(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := generateInvoiceForm{}

		err := c.BodyParser(&form)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{
				"the JSON structure is not correct",
			})
		}

		clientID := uint(form.Header.ClientID)

		_, err = s.User.Find(clientID)
		if errors.Is(err, domain.ErrUserNotFound) {
			// TO-DO
			//c.Logger().Error("User not found to generate invoice")

			return c.Status(http.StatusNotFound).JSON(respError{err.Error()}) // 404
		}

		if err != nil {
			// TO-DO
			//c.Logger().Error(err)

			return c.Status(http.StatusBadRequest).JSON(respError{err.Error()})
		}

		assemble := func(i invoiceItemForm) *domain.InvoiceItem {
			return &domain.InvoiceItem{
				ProductID: i.ProductID,
			}
		}

		items := make(domain.ItemList, 0, len(form.Items))
		for _, item := range form.Items {

			_, err := s.Store.Find(item.ProductID)

			if errors.Is(err, domain.ErrProductNotFound) {
				// TO-DO
				//c.Logger().Error("Product not found to add the invoice")

				return c.Status(http.StatusNotFound).JSON(respError{
					fmt.Sprintf("%s with id %d", domain.ErrProductNotFound, item.ProductID),
				}) // 404
			}

			if err != nil {
				// TO-DO
				//c.Logger().Error(err)

				return c.Status(http.StatusBadRequest).JSON(respError{err.Error()})
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
			// TO-DO
			//c.Logger().Error(err)

			return c.Status(http.StatusInternalServerError).JSON(respError{err.Error()})
		}

		// TO-DO
		//c.Logger().Infof("Invoice ID %d generated", invoice.Header.ID)

		return c.Status(http.StatusCreated).JSON(respOkData{
			Msg: "invoice generated",
			Data: generateInvoiceForm{
				Header: form.Header,
				Items:  form.Items,
			},
		})
	}
}

// generateInvoiceForm models of fields to request to generate an invoice.
type generateInvoiceForm struct {
	Header invoiceHeaderForm `json:"header"`
	Items  []invoiceItemForm `json:"items"`
}

// invoiceItemForm represents a form to generate invoice item as product.
type invoiceItemForm struct {
	ProductID int `json:"productId"`
}

type invoiceHeaderForm struct {
	ClientID int `json:"clientId"`
}

// InvoiceReportDTO represent a view of a invoice.
type InvoiceReportDTO struct {
	Header invoiceHeaderForm  `json:"header"`
	Items  []*invoiceItemForm `json:"items"`
}

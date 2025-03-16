package http

import (
	"errors"
	"fmt"
	"net/http"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/app"
	"github.com/adrianolmedo/genesis/logger"

	"github.com/gofiber/fiber/v2"
)

// generateInvoice godoc
//
//	@Summary		Generate invoice
//	@Description	Generate invoice of products
//	@Tags			billing
//	@Accept			json
//	@Produce		json
//	@Failure		400					{object}	errorResp
//	@Failure		404					{object}	errorResp
//	@Failure		500					{object}	errorResp
//	@Success		201					{object}	successResp{data=generateInvoiceForm}
//	@Param			generateInvoiceForm	body		generateInvoiceForm	true	"application/json"
//	@Router			/invoices [post]
func generateInvoice(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := generateInvoiceForm{}

		err := c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		clientID := uint(form.Header.ClientID)

		_, err = s.User.Find(clientID)
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("generating invoice", fmt.Sprintf("user ID %d not found to generate invoice", clientID))

			return errorJSON(c, http.StatusNotFound, respDetails{
				Code:    "002",
				Message: err.Error(),
			})
		}

		if err != nil {
			logger.Error("generating invoice", err.Error())

			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: err.Error(),
			})
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
				logger.Debug("generating invoice", fmt.Sprintf("product ID %d not found to add the invoice", item.ProductID))

				return errorJSON(c, http.StatusNotFound, respDetails{
					Code:    "002",
					Message: fmt.Sprintf("%s with id %d", domain.ErrProductNotFound, item.ProductID),
				})
			}

			if err != nil {
				logger.Error("generating invoice", err.Error())

				return errorJSON(c, http.StatusBadRequest, respDetails{
					Code:    "002",
					Message: err.Error(),
				})
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
			logger.Error("generating invoice", err.Error())

			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "002",
				Message: err.Error(),
			})
		}

		logger.Info("generating invoice", fmt.Sprintf("invoice ID %d generated", invoice.Header.ID))

		return successJSON(c, http.StatusCreated, respDetails{
			Message: "Invoice generated",
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

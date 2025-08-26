package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/adrianolmedo/genesis/billing"
	"github.com/adrianolmedo/genesis/compose"
	"github.com/adrianolmedo/genesis/logger"
	"github.com/adrianolmedo/genesis/store"
	"github.com/adrianolmedo/genesis/user"

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
//	@Success		201					{object}	resp{data=generateInvoiceReq}
//	@Param			generateInvoiceReq	body		generateInvoiceReq	true	"application/json"
//	@Router			/invoices [post]
func generateInvoice(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := generateInvoiceReq{}

		err := c.BodyParser(&req)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		clientID := req.Header.ClientID

		_, err = svcs.User.Find(ctx, clientID)
		if errors.Is(err, user.ErrNotFound) {
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

		// assmble is like a mapper to convert invoiceItemReq to billing.InvoiceItem.
		assemble := func(i invoiceItemReq) *billing.InvoiceItem {
			return &billing.InvoiceItem{
				ProductID: i.ProductID,
			}
		}

		items := make(billing.ItemList, 0, len(req.Items))
		for _, item := range req.Items {

			_, err := svcs.Store.Find(ctx, item.ProductID)

			if errors.Is(err, store.ErrProductNotFound) {
				logger.Debug("generating invoice", fmt.Sprintf("product ID %d not found to add the invoice", item.ProductID))

				return errorJSON(c, http.StatusNotFound, respDetails{
					Code:    "002",
					Message: fmt.Sprintf("%s with id %d", store.ErrProductNotFound, item.ProductID),
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

		invoice := &billing.Invoice{
			Header: &billing.InvoiceHeader{
				ClientID: clientID,
			},
			Items: items,
		}

		err = svcs.Billing.Generate(ctx, invoice)
		if err != nil {
			logger.Error("generating invoice", err.Error())

			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "002",
				Message: err.Error(),
			})
		}

		logger.Info("generating invoice", fmt.Sprintf("invoice ID %d generated", invoice.Header.ID))

		return respJSON(c, http.StatusCreated, respDetails{
			Message: "Invoice generated",
			Data: generateInvoiceReq{
				Header: req.Header,
				Items:  req.Items,
			},
		})
	}
}

// generateInvoiceReq models of fields to request to generate an invoice.
type generateInvoiceReq struct {
	Header invoiceHeaderReq `json:"header"`
	Items  []invoiceItemReq `json:"items"`
}

// invoiceItemReq represents a Command to generate invoice item as product.
type invoiceItemReq struct {
	ProductID int64 `json:"productId"`
}

type invoiceHeaderReq struct {
	ClientID int64 `json:"clientId"`
}

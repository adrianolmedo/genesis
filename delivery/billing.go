package delivery

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/adrianolmedo/go-restapi/domain"

	"github.com/gofiber/fiber/v2"
)

func generateInvoice(s *services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := domain.GenerateInvoiceForm{}

		err := c.BodyParser(&form)
		if err != nil {
			resp := respJSON(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		clientID := form.Header.ClientID

		_, err = s.User.Find(clientID)
		if errors.Is(err, domain.ErrUserNotFound) {
			// TO-DO
			//c.Logger().Error("User not found to generate invoice")

			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNotFound).JSON(resp) // 404
		}

		if err != nil {
			// TO-DO
			//c.Logger().Error(err)

			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
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
				// TO-DO
				//c.Logger().Error("Product not found to add the invoice")

				resp := respJSON(msgError, fmt.Sprintf("%s with id %d", domain.ErrProductNotFound, item.ProductID), nil)
				return c.Status(http.StatusNotFound).JSON(resp) // 404
			}

			if err != nil {
				// TO-DO
				//c.Logger().Error(err)

				resp := respJSON(msgError, err.Error(), nil)
				return c.Status(http.StatusBadRequest).JSON(resp)
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

			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO
		//c.Logger().Infof("Invoice ID %d generated", invoice.Header.ID)

		resp := respJSON(msgOK, "invoice generated", domain.GenerateInvoiceForm{
			Header: form.Header,
			Items:  form.Items,
		})
		return c.Status(http.StatusCreated).JSON(resp)
	}
}

package gorestapi

import (
	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/postgres"
)

type BillingService struct {
	repo postgres.Invoice
}

func (s BillingService) Generate(invoice *domain.Invoice) error {
	err := generateInvoice(invoice)
	if err != nil {
		return err
	}

	return s.repo.Create(invoice)
}

func generateInvoice(invoice *domain.Invoice) error {
	if invoice.Items == nil || len(invoice.Items) == 0 {
		return domain.ErrItemListCantBeEmpty
	}

	return nil
}

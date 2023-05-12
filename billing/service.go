package billing

import (
	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/postgres"
)

type Service struct {
	repo postgres.InvoiceRepository
}

func NewBillingService(repo postgres.InvoiceRepository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) Generate(invoice *domain.Invoice) error {
	err := generateInvoiceService(invoice)
	if err != nil {
		return err
	}

	return s.repo.Create(invoice)
}

func generateInvoiceService(invoice *domain.Invoice) error {
	if invoice.Items == nil || len(invoice.Items) == 0 {
		return domain.ErrItemListCantBeEmpty
	}

	return nil
}

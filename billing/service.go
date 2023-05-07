package billing

import (
	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/storage"
)

type Service struct {
	repo storage.InvoiceRepository
}

func NewBillingService(repo storage.InvoiceRepository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) Generate(invoice *domain.Invoice) error {
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

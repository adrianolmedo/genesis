package service

import (
	"github.com/adrianolmedo/go-restapi/internal/domain"
	"github.com/adrianolmedo/go-restapi/internal/storage"
)

type BillingService interface {
	Generate(*domain.Invoice) error
}

type billingService struct {
	repo storage.InvoiceRepository
}

func NewBillingService(repo storage.InvoiceRepository) BillingService {
	return &billingService{repo}
}

func (bs billingService) Generate(invoice *domain.Invoice) error {
	if invoice.Items == nil || len(invoice.Items) == 0 {
		return domain.ErrItemListCantBeEmpty
	}

	return bs.repo.Create(invoice)
}

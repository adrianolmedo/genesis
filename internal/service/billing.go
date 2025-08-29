package service

import (
	"github.com/adrianolmedo/genesis/internal/domain"
	"github.com/adrianolmedo/genesis/internal/storage"
)

type BillingService interface {
	Generate(*domain.Invoice) error
}

type billingService struct {
	repo storage.InvoiceRepository
}

func NewBillingService(repo storage.InvoiceRepository) *billingService {
	return &billingService{repo}
}

func (bs billingService) Generate(invoice *domain.Invoice) error {
	if invoice.Items == nil || len(invoice.Items) == 0 {
		return domain.ErrItemListCantBeEmpty
	}

	return bs.repo.Create(invoice)
}

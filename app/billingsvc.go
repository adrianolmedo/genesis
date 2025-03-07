package app

import (
	domain "github.com/adrianolmedo/genesis"
	strg "github.com/adrianolmedo/genesis/pgsql/pq"
)

type billingService struct {
	repo strg.Invoice
}

func (s billingService) Generate(inv *domain.Invoice) error {
	err := generateInvoice(inv)
	if err != nil {
		return err
	}

	return s.repo.Create(inv)
}

func generateInvoice(inv *domain.Invoice) error {
	if inv.Items == nil || inv.Items.IsEmpty() {
		return domain.ErrItemListCantBeEmpty
	}

	return nil
}

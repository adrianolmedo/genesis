package app

import (
	"context"

	domain "github.com/adrianolmedo/genesis"
	storage "github.com/adrianolmedo/genesis/pgsql/sqlc"
)

type billingService struct {
	repo *storage.Invoice
}

func (s billingService) Generate(ctx context.Context, inv *domain.Invoice) error {
	err := generateInvoice(inv)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, inv)
}

func generateInvoice(inv *domain.Invoice) error {
	if inv.Items == nil || inv.Items.IsEmpty() {
		return domain.ErrItemListCantBeEmpty
	}

	return nil
}

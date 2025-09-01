package billing

import "context"

type Service struct {
	repo *Repo
}

func NewService(r *Repo) *Service {
	return &Service{repo: r}
}

func (s Service) Generate(ctx context.Context, inv *Invoice) error {
	err := generateInvoice(inv)
	if err != nil {
		return err
	}

	return s.repo.CreateInvoice(ctx, inv)
}

func generateInvoice(inv *Invoice) error {
	if inv.Items.IsEmpty() {
		return ErrItemListCantBeEmpty
	}

	return nil
}

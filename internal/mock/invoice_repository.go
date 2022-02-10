package mock

import (
	"errors"

	"github.com/adrianolmedo/go-restapi/internal/domain"
)

type InvoiceRepositoryOk struct{}

func (InvoiceRepositoryOk) Create(*domain.Invoice) error {
	return nil
}

type InvoiceRepositoryError struct{}

func (InvoiceRepositoryError) Create(*domain.Invoice) error {
	return errors.New("mock error")
}

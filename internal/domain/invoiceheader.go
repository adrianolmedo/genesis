package domain

import (
	"errors"
	"time"
)

var ErrInvoiceHeaderNotFound = errors.New("invoice header not found")

type InvoiceHeader struct {
	ID       int64
	ClientID int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

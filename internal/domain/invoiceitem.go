package domain

import (
	"errors"
	"time"
)

var ErrInvoiceItemNotFound = errors.New("invoice item not found")

type InvoiceItem struct {
	ID              int64
	InvoiceHeaderID int64
	ProductID       int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ItemList []*InvoiceItem

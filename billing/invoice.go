package billing

import (
	"errors"
	"time"
)

// ErrInvoiceHeaderNotFound a header related with a invoice it's not found.
var ErrInvoiceHeaderNotFound = errors.New("invoice header not found")

// ErrInvoiceItemNotFound a item related with a invoice it's not found.
var ErrInvoiceItemNotFound = errors.New("invoice item not found")

// ErrItemListCantBeEmpty indicate that an invoice must have items.
var ErrItemListCantBeEmpty = errors.New("item list can't be empty")

// Invoice aggregate, represents an invoice with its header and items.
type Invoice struct {
	Header *InvoiceHeader
	Items  ItemList
}

// InvoiceHeader entity model.
type InvoiceHeader struct {
	ID       int64
	UUID     string
	ClientID int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

// InvoiceItem entity model.
type InvoiceItem struct {
	ID              int64
	InvoiceHeaderID int64
	ProductID       int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ItemList collection of invoice items.
type ItemList []InvoiceItem

// IsEmpty return true if is empty.
func (il ItemList) IsEmpty() bool {
	return len(il) == 0
}

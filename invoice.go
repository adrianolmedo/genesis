package genesis

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

// Invoice model.
type Invoice struct {
	Header *InvoiceHeader
	Items  ItemList
}

// InvoiceHeader model.
type InvoiceHeader struct {
	ID       int
	UUID     string
	ClientID uint

	CreatedAt time.Time
	UpdatedAt time.Time
}

// InvoiceItem model.
type InvoiceItem struct {
	ID              int
	InvoiceHeaderID int
	ProductID       int

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ItemList collection of invoice items.
type ItemList []*InvoiceItem

// IsEmpty return true if is empty.
func (il ItemList) IsEmpty() bool {
	return len(il) == 0
}

// InvoiceItemForm represents a form to generate invoice item as product.
// TODO: Pass DTO to http/ package.
type InvoiceItemForm struct {
	ProductID int `json:"productId"`
}

// GenerateInvoiceForm models of fields to request to generate an invoice.
// TODO: Pass DTO to http/ package.
type GenerateInvoiceForm struct {
	Header invoiceHeaderForm `json:"header"`
	Items  []InvoiceItemForm `json:"items"`
}

type invoiceHeaderForm struct {
	ClientID int `json:"clientId"`
}

// InvoiceReportDTO represent a view of a invoice.
// TODO: Pass DTO to http/ package.
type InvoiceReportDTO struct {
	Header invoiceHeaderForm  `json:"header"`
	Items  []*InvoiceItemForm `json:"items"`
}

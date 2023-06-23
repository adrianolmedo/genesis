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
	ID       int64
	ClientID int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

// InvoiceItem model.
type InvoiceItem struct {
	ID              int64
	InvoiceHeaderID int64
	ProductID       int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ItemList collection of invoice items.
type ItemList []*InvoiceItem

// InvoiceItemForm represents a form to generate invoice item as product.
type InvoiceItemForm struct {
	ProductID int64 `json:"product_id"`
}

// GenerateInvoiceForm models of fields to request to generate an invoice.
type GenerateInvoiceForm struct {
	Header invoiceHeaderForm `json:"header"`
	Items  []InvoiceItemForm `json:"items"`
}

type invoiceHeaderForm struct {
	ClientID int64 `json:"client_id"`
}

// InvoiceReportDTO represent a view of a invoice.
type InvoiceReportDTO struct {
	Header invoiceHeaderForm  `json:"header"`
	Items  []*InvoiceItemForm `json:"items"`
}

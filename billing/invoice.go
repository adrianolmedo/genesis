package billing

import (
	"errors"
	"time"
)

var ErrInvoiceHeaderNotFound = errors.New("invoice header not found")
var ErrInvoiceItemNotFound = errors.New("invoice item not found")
var ErrItemListCantBeEmpty = errors.New("item list can't be empty")

type Invoice struct {
	Header *InvoiceHeader
	Items  ItemList
}

type InvoiceHeader struct {
	ID       int64
	ClientID int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type InvoiceItem struct {
	ID              int64
	InvoiceHeaderID int64
	ProductID       int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ItemList []*InvoiceItem

type InvoiceHeaderForm struct {
	ClientID int64 `json:"client_id"`
}

type InvoiceItemForm struct {
	ProductID int64 `json:"product_id"`
}

type ItemListForm []*InvoiceItemForm

type GenerateInvoiceForm struct {
	Header InvoiceHeaderForm `json:"header"`
	Items  ItemListForm      `json:"items"`
}

type InvoiceReportDTO struct {
	Header InvoiceHeaderForm `json:"header"`
	Items  ItemListForm      `json:"items"`
}

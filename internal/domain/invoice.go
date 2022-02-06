package domain

type Invoice struct {
	Header *InvoiceHeader
	Items  ItemList
}

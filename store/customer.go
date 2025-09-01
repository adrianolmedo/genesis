package store

import (
	"errors"

	"github.com/adrianolmedo/genesis"
)

var ErrCustomerNotFound = errors.New("customer not found")

// Customer domain model.
type Customer struct {
	ID        int64
	UUID      string
	FirstName string
	LastName  string
	Password  string
	Email     string

	genesis.AuditFields
}

// Validate return error if certain fields there are empty.
func (c Customer) Validate() error {
	if c.FirstName == "" || c.Email == "" {
		return errors.New("first name, email can't be empty")
	}
	return nil
}

// Customers collection of Customer.
type Customers []Customer

// IsEmpty return true if is empty.
func (cs Customers) IsEmpty() bool {
	return len(cs) == 0
}

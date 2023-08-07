package genesis

import (
	"errors"
	"time"
)

var ErrCustomerNotFound = errors.New("customer not found")

// Customer model.
type Customer struct {
	ID        int
	UUID      string
	FirstName string
	LastName  string
	Email     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// Validate return error if certain fields there are empty.
func (c Customer) Validate() error {
	if c.FirstName == "" || c.Email == "" {
		return errors.New("first name, email can't be empty")
	}
	return nil
}

// Customers collection of Customer.
type Customers []*Customer

// IsEmpty return true if is empty.
func (cs Customers) IsEmpty() bool {
	return len(cs) == 0
}

// CreateCustomerForm subset of fields to request to create a Customer.
// TODO: Pass DTO to http/ package.
type CreateCustomerForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// CustomerProfileDTO subset of Customer fields.
// TODO: Pass DTO to http/ package.
type CustomerProfileDTO struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

package genesis

import (
	"errors"
	"time"
)

// Customer model.
type Customer struct {
	ID        int64
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
type Customers []Customer

// IsEmpty return true if is empty.
func (cs Customers) IsEmpty() bool {
	return len(cs) == 0
}

// CreateCustomerForm subset of fields to request to create a Customer.
type CreateCustomerForm struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// CustomerProfileDTO subset of Customer fields.
type CustomerProfileDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// CustomerList collection of Customer presentation.
type CustomerList []CustomerProfileDTO

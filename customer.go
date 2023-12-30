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
	Password  string
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

package genesis

import (
	"errors"
	"time"

	"github.com/pborman/uuid"
)

var ErrUserCantBeEmpty = errors.New("the user fields can't be empty")
var ErrUserNotFound = errors.New("user not found")

// User domain model.
type User struct {
	ID        int64
	UUID      string
	FirstName string
	LastName  string
	Email     string
	Password  string

	AuditFields
}

// Validate return error if certain fields there are empty.
func (u User) Validate() error {
	if u.FirstName == "" || u.Email == "" || u.Password == "" {
		return errors.New("first name, email or password can't be empty")
	}

	return nil
}

// Users a collection of User.
type Users []*User

// IsEmpty return true if is empty.
func (us Users) IsEmpty() bool {
	return len(us) == 0
}

// NextUUID generates a new UUID.
func NextUUID() string {
	return uuid.New()
}

type AuditFields struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

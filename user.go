package genesis

import (
	"errors"
	"time"

	"github.com/pborman/uuid"
)

var ErrUserCantBeEmpty = errors.New("the user fields can't be empty")
var ErrUserNotFound = errors.New("user not found")

// User model.
type User struct {
	ID        int
	UUID      string
	FirstName string
	LastName  string
	Email     string
	Password  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
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

// UserSignUpForm subset of User fields to create account.
// TODO: Pass DTO to http/ package.
type UserSignUpForm struct {
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Email     string `json:"email" example:"johndoe@aol.com"`
	Password  string `json:"password" example:"1234567b"`
}

// UserLoginForm subset of user fields to request login.
// TODO: Pass DTO to http/ package.
type UserLoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserUpdateForm subset of fields to request to update a User.
// TODO: Pass DTO to http/ package.
type UserUpdateForm struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// UserProfileDTO subset of User fields.
// TODO: Pass DTO to http/ package.
type UserProfileDTO struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

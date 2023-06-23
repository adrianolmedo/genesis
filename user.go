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
	ID        int64
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

// NextUUID generates a new UUID.
func NextUUID() string {
	return string(uuid.New())
}

// UserSignUpForm subset of User fields to create account.
type UserSignUpForm struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// UserLoginForm subset of user fields to request login.
type UserLoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserUpdateForm subset of fields to request to update a User.
type UserUpdateForm struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// UserProfileDTO subset of User fields .
type UserProfileDTO struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// UserList collection of User presentation.
type UserList []UserProfileDTO

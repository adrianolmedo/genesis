package user

import (
	"errors"
	"time"
)

var ErrUserCantBeEmpty = errors.New("the user fields can't be empty")
var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Password  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// CheckEmptyFields return error if FirstName, Email or Password there are empty.
func (u User) CheckEmptyFields() error {
	if u.FirstName == "" || u.Email == "" || u.Password == "" {
		return errors.New("first name, email or password can't be empty")
	}
	return nil
}

type UserProfileDTO struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
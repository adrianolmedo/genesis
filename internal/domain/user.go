package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

var ErrUserCantBeEmpty = errors.New("the user fields can't be empty")
var ErrUserNotFound = errors.New("user not found")

// User is the central struct in the domain model.
type User struct {
	ID        int64
	UUID      UserID
	FirstName string
	LastName  string
	Email     string
	Password  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// CheckEmptyFields return error if FirstName, Email or Password there are empty.
// TO-DO: Change to boolean function.
func (u User) CheckEmptyFields() error {
	if u.FirstName == "" || u.Email == "" || u.Password == "" {
		return errors.New("first name, email or password can't be empty")
	}
	return nil
}

// UserID uniquely identifies a particular user.
type UserID string

// NextUserID generates a new UUID.
func NextUserID() UserID {
	return UserID(strings.Split(strings.ToUpper(uuid.New()), "-")[0])
}

type UserSignUpForm struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserLoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateForm struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserProfileDTO struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UsersList []UserProfileDTO

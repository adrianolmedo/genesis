package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/pborman/uuid"
)

// ErrResourceCantBeEmpty = "the resource can't be empty"
var ErrResourceCantBeEmpty = errors.New("the resource can't be empty")

// ErrResourceDoesNotExist = "resource does not exist".
var ErrResourceDoesNotExist = errors.New("resource does not exist")

// User is the central struct in the domain model.
type User struct {
	ID        int64
	UUID      UserID
	FirstName string
	LastName  string
	Email     string
	Password  string
}

// checkEmptyFields return error if FirstName, Email or Password there are empty.
func (u User) checkEmptyFields() error {
	if u.FirstName == "" || u.Email == "" || u.Password == "" {
		return errors.New("first name, email or password can't be empty")
	}
	return nil
}

func (u User) validateEmail() error {
	validEmail, err := regexp.MatchString(`^([a-zA-Z0-9])+([a-zA-Z0-9\._-])*@([a-zA-Z0-9_-])+([a-zA-Z0-9\._-]+)+$`, u.Email)
	if err != nil {
		return fmt.Errorf("email pattern: %v", err)
	}

	if !validEmail {
		return errors.New("email not valid")
	}
	return nil
}

// UserID uniquely identifies a particular user.
type UserID string

// NextUserID generates a new UUID.
func NextUserID() UserID {
	return UserID(strings.Split(strings.ToUpper(uuid.New()), "-")[0])
}

type ProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateRequest struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type ProfileResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type ListResponse []ProfileResponse

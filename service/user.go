package service

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"go-rest-practice/model"
)

// UserDAO interface is InterfaceDAO in Pattern DAO.
// http://chuwiki.chuidiang.org/index.php?title=Patr%C3%B3n_DAO
type UserDAO interface {
	Create(*model.User) error
	ByID(id int64) (*model.User, error)
	Update(*model.User) error
	All() (model.Users, error)
	Delete(id int64) error
}

// UserService para desacoplar la comunicación entre UserDAO y storage.User.
type User struct {
	storage UserDAO
}

func NewUser(s UserDAO) *User {
	return &User{
		storage: s,
	}
}

// Create is business logic for create a User.
func (s User) Create(user *model.User) error {
	err := validateUser(user)
	if err != nil {
		return err
	}

	user.CreatedAt = time.Now()
	return s.storage.Create(user)
}

// ByID is business logic for get a User by it's ID.
func (s User) ByID(id int64) (*model.User, error) {
	return s.storage.ByID(id)
}

// Update is business logic for update a User.
func (s User) Update(user *model.User) error {
	err := validateUser(user)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	return s.storage.Update(user)
}

func (s User) All() (model.Users, error) {
	return s.storage.All()
}

// Delete is business logic for delete User by it's ID.
func (s User) Delete(id int64) error {
	return s.storage.Delete(id)
}

// validateUser it ensures that User fields like
// FirtsName, LastName, Email and Password can't be empty.
func validateUser(user *model.User) error {
	if user == nil {
		return model.ErrResourceCantBeEmpty
	}

	if user.FirstName == "" || user.Email == "" || user.Password == "" {
		return errors.New("first name, email or password can't be empty")
	}
	return validateEmail(user.Email)
}

// validateEmail check if a string is a valid email.
func validateEmail(email string) error {
	validEmail, err := regexp.MatchString(`^([a-zA-Z0-9])+([a-zA-Z0-9\._-])*@([a-zA-Z0-9_-])+([a-zA-Z0-9\._-]+)+$`, email)
	if err != nil {
		return fmt.Errorf("email pattern: %v", err)
	}

	if !validEmail {
		return errors.New("email not valid")
	}
	return nil
}

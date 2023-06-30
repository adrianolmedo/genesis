package app

import (
	"errors"
	"fmt"
	"regexp"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/postgres"
)

type userService struct {
	repo postgres.User
}

func (s userService) Login(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	return s.repo.ByLogin(email, password)
}

// SignUp to register a User.
func (s userService) SignUp(u *domain.User) error {
	err := signUp(u)
	if err != nil {
		return err
	}

	return s.repo.Create(u)
}

// signUp applicaction logic for regitser a User. Has been split into
// a smaller function for unit testing purposes, and it should do so for
// the other methods of the Service.
func signUp(u *domain.User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	err = validateEmail(u.Email)
	if err != nil {
		return err
	}

	return nil
}

// Find a User by its ID.
func (s userService) Find(id int) (*domain.User, error) {
	if id == 0 {
		return &domain.User{}, domain.ErrUserNotFound
	}

	return s.repo.ByID(id)
}

// Update application logic for update a User.
func (s userService) Update(u domain.User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	err = validateEmail(u.Email)
	if err != nil {
		return err
	}

	return s.repo.Update(u)
}

// List get list of users.
func (s userService) List() (domain.Users, error) {
	return s.repo.All()
}

// Remove delete User by its ID.
func (s userService) Remove(id int64) error {
	if id == 0 {
		return domain.ErrUserNotFound
	}

	return s.repo.Delete(id)
}

// validateEmail helper to check email pattern.
func validateEmail(email string) error {
	validEmail, err := regexp.MatchString(`^([a-zA-Z0-9])+([a-zA-Z0-9\._-])*@([a-zA-Z0-9_-])+([a-zA-Z0-9\._-]+)+$`, email)
	if err != nil {
		return fmt.Errorf("email pattern: %v", err)
	}

	if !validEmail {
		return errors.New("invalid email")
	}

	return nil
}

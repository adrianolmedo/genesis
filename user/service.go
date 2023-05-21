package user

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/adrianolmedo/go-restapi/domain"
	"github.com/adrianolmedo/go-restapi/postgres"
)

type Service struct {
	repo postgres.User
}

func NewService(repo postgres.User) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) Login(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	return s.repo.ByLogin(email, password)
}

// SignUp to register a User.
func (s Service) SignUp(u *domain.User) error {
	err := signUp(u)
	if err != nil {
		return err
	}

	return s.repo.Create(u)
}

// signUp business logic for regitser a User. Has been split into
// a smaller function for unit testing purposes, and it should do so for
// the other methods of the Service.
func signUp(u *domain.User) error {
	err := u.CheckEmptyFields()
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
func (s Service) Find(id int64) (*domain.User, error) {
	if id == 0 {
		return &domain.User{}, domain.ErrUserNotFound
	}

	return s.repo.ByID(id)
}

// Update business logic for update a User.
func (s Service) Update(u domain.User) error {
	err := u.CheckEmptyFields()
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
func (s Service) List() ([]*domain.User, error) {
	return s.repo.All()
}

// Remove delete User by its ID.
func (s Service) Remove(id int64) error {
	if id == 0 {
		return domain.ErrUserNotFound
	}

	return s.repo.Delete(id)
}

// helpers

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

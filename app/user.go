package app

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/pgsql"
	storage "github.com/adrianolmedo/genesis/pgsql/sqlc"
)

type userService struct {
	repo *storage.User
}

func (s userService) Login(ctx context.Context, email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	return s.repo.ByLogin(ctx, email, password)
}

// SignUp to register a User.
func (s userService) SignUp(ctx context.Context, u *domain.User) error {
	err := signUp(u)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, u)
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
func (s userService) Find(ctx context.Context, id int64) (*domain.User, error) {
	if id == 0 {
		return &domain.User{}, domain.ErrUserNotFound
	}

	return s.repo.ByID(ctx, id)
}

// Update application logic for update a User.
func (s userService) Update(ctx context.Context, u domain.User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	err = validateEmail(u.Email)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, u)
}

// List get list of users.
func (s userService) List(ctx context.Context, p *pgsql.Pager) (pgsql.PagerResults, error) {
	return s.repo.List(ctx, p)
}

// Remove delete User by its ID.
func (s userService) Remove(ctx context.Context, id int64) error {
	if id == 0 {
		return domain.ErrUserNotFound
	}

	return s.repo.Delete(ctx, id)
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

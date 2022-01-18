package service

import (
	"time"

	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage"
)

// UserService (before UserDAO) interface is InterfaceDAO in Pattern DAO.
// http://chuwiki.chuidiang.org/index.php?title=Patr%C3%B3n_DAO
type UserService interface {
	SignUp(*domain.User) error
	Find(id int64) (*domain.User, error)
	Update(domain.User) error
	List() ([]*domain.User, error)
	Remove(id int64) error
}

type userService struct {
	repo storage.UserRepository
}

func NewUserService(repo storage.UserRepository) UserService {
	return &userService{repo}
}

// SignUp business logic to register a User.
func (us userService) SignUp(user *domain.User) error {
	err := user.CheckEmptyFields()
	if err != nil {
		return err
	}

	err = validateEmail(user.Email)
	if err != nil {
		return err
	}

	user.UUID = domain.NextUserID()
	user.CreatedAt = time.Now()

	return us.repo.Create(user)
}

// Find a User by its ID.
func (us userService) Find(id int64) (*domain.User, error) {
	return us.repo.ByID(id)
}

// Update business logic for update a User.
func (us userService) Update(user domain.User) error {
	err := user.CheckEmptyFields()
	if err != nil {
		return err
	}

	err = validateEmail(user.Email)
	if err != nil {
		return err
	}

	return us.repo.Update(user)
}

// List get list of users.
func (us userService) List() ([]*domain.User, error) {
	return us.repo.All()
}

// Remove delete User by its ID.
func (us userService) Remove(id int64) error {
	return us.repo.Delete(id)
}

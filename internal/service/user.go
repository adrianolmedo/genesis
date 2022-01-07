package service

import (
	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/repository"
)

// UserService (before UserDAO) interface is InterfaceDAO in Pattern DAO.
// http://chuwiki.chuidiang.org/index.php?title=Patr%C3%B3n_DAO
type UserService interface {
	SignUp(domain.User) error
	Find(id int64) (*domain.User, error)
	Update(domain.User) error
	List() ([]*domain.User, error)
	Remove(id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

// SignUp is business logic to register a User.
func (us userService) SignUp(user domain.User) error {
	err := user.CheckEmptyFields()
	if err != nil {
		return err
	}

	err = validateEmail(user.Email)
	if err != nil {
		return err
	}

	user.UUID = domain.NextUserID()
	return us.repo.Create(user)
}

// ByID is business logic for get a User by its ID.
func (us userService) Find(id int64) (*domain.User, error) {
	return us.repo.ByID(id)
}

// Update is business logic for update a User.
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

func (us userService) List() ([]*domain.User, error) {
	return us.repo.All()
}

// Remove is business logic for delete User by its ID.
func (us userService) Remove(id int64) error {
	return us.repo.Delete(id)
}

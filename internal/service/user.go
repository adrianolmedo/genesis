package service

import (
	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage"
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
	repository storage.UserRepository
}

func NewUserService(r storage.UserRepository) UserService {
	return &userService{r}
}

// SignUp is business logic to register a User.
func (s userService) SignUp(user domain.User) error {
	err := user.CheckEmptyFields()
	if err != nil {
		return err
	}

	err = validateEmail(user.Email)
	if err != nil {
		return err
	}

	user.UUID = domain.NextUserID()
	return s.repository.Create(user)
}

// ByID is business logic for get a User by its ID.
func (s userService) Find(id int64) (*domain.User, error) {
	return s.repository.ByID(id)
}

// Update is business logic for update a User.
func (s userService) Update(user domain.User) error {
	err := user.CheckEmptyFields()
	if err != nil {
		return err
	}

	err = validateEmail(user.Email)
	if err != nil {
		return err
	}

	return s.repository.Update(user)
}

func (s userService) List() ([]*domain.User, error) {
	return s.repository.All()
}

// Remove is business logic for delete User by its ID.
func (s userService) Remove(id int64) error {
	return s.repository.Delete(id)
}

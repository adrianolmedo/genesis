package user

import (
	"errors"
	"fmt"
	"regexp"
)

// Service (before UserDAO) interface is InterfaceDAO in Pattern DAO.
// http://chuwiki.chuidiang.org/index.php?title=Patr%C3%B3n_DAO
type Service interface {
	Register(*Request) error
	ByID(id int64) (*Response, error)
	Update(*Request) error
	All() ([]*Response, error)
	Delete(id int64) error
}

// Repository para desacoplar el paquete `repository` de persistencia
// de datos entre postgres o mysql.
//
// Si no hubiera persistencia de datos, esta interface no debería existir,
// y en su lugar se debería importar directamente el modelo, domain u objeto de dominio/valor
// como un campo en la estructura `service`.
type Repository interface {
	Create(*Request) error
	ByID(id int64) (*Response, error)
	Update(*Request) error
	All() ([]*Response, error)
	Delete(id int64) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

// Create is business logic for create a User.
func (s service) Register(user *Request) error {
	err := validateUser(user)
	if err != nil {
		return err
	}

	//user.CreatedAt = time.Now()
	return s.repository.Create(user)
}

// ByID is business logic for get a User by it's ID.
func (s service) ByID(id int64) (*Response, error) {
	return s.repository.ByID(id)
}

// Update is business logic for update a User.
func (s service) Update(user *Request) error {
	err := validateUser(user)
	if err != nil {
		return err
	}

	//user.UpdatedAt = time.Now()
	return s.repository.Update(user)
}

func (s service) All() ([]*Response, error) {
	return s.repository.All()
}

// Delete is business logic for delete User by it's ID.
func (s service) Delete(id int64) error {
	return s.repository.Delete(id)
}

// validateUser it ensures that User fields like
// FirtsName, LastName, Email and Password can't be empty.
func validateUser(user *Request) error {
	if user == nil {
		return ErrResourceCantBeEmpty
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

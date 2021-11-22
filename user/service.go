package user

// Service (before UserDAO) interface is InterfaceDAO in Pattern DAO.
// http://chuwiki.chuidiang.org/index.php?title=Patr%C3%B3n_DAO
type Service interface {
	Register(*User) error
	ByID(id int64) (*User, error)
	Update(*User) error
	All() ([]*User, error)
	Delete(id int64) error
}

// Repository to uncouple persistence `repository` package
// data between postgres or mysql.
//
// If there is no data persistence, this interface should not exist,
// and instead, the model, domain or domain/value object should be imported
// from `repository` directly as field in the `service` struct.
type Repository interface {
	Create(*User) error
	ByID(id int64) (*User, error)
	Update(*User) error
	All() ([]*User, error)
	Delete(id int64) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

// Create is business logic for create a User.
func (s service) Register(u *User) error {
	err := u.checkEmptyFields()
	if err != nil {
		return err
	}

	err = u.validateEmail()
	if err != nil {
		return err
	}

	u.UUID = NextUserID()
	return s.repository.Create(u)
}

// ByID is business logic for get a User by its ID.
func (s service) ByID(id int64) (*User, error) {
	return s.repository.ByID(id)
}

// Update is business logic for update a User.
func (s service) Update(u *User) error {
	err := u.checkEmptyFields()
	if err != nil {
		return err
	}

	err = u.validateEmail()
	if err != nil {
		return err
	}

	return s.repository.Update(u)
}

func (s service) All() ([]*User, error) {
	return s.repository.All()
}

// Delete is business logic for delete User by its ID.
func (s service) Delete(id int64) error {
	return s.repository.Delete(id)
}

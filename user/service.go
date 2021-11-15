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

// Repository para desacoplar el paquete `repository` de persistencia
// de datos entre postgres o mysql.
//
// Si no hubiera persistencia de datos, esta interface no debería existir,
// y en su lugar se debería importar directamente el modelo, domain u objeto de dominio/valor
// como un campo en la estructura `service`.
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
func (s service) Register(user *User) error {
	err := user.validateNoEmptyFields()
	if err != nil {
		return err
	}

	err = user.validateEmail()
	if err != nil {
		return err
	}

	//user.CreatedAt = time.Now()
	return s.repository.Create(user)
}

// ByID is business logic for get a User by it's ID.
func (s service) ByID(id int64) (*User, error) {
	return s.repository.ByID(id)
}

// Update is business logic for update a User.
func (s service) Update(user *User) error {
	err := user.validateNoEmptyFields()
	if err != nil {
		return err
	}

	err = user.validateEmail()
	if err != nil {
		return err
	}

	//user.UpdatedAt = time.Now()
	return s.repository.Update(user)
}

func (s service) All() ([]*User, error) {
	return s.repository.All()
}

// Delete is business logic for delete User by it's ID.
func (s service) Delete(id int64) error {
	return s.repository.Delete(id)
}

package user

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

// SignUp to register a User.
func (s Service) SignUp(user *User) error {
	err := signUpService(user)
	if err != nil {
		return err
	}

	return s.repo.Create(user)
}

// signUpService business logic for regitser a User. Has been split into
// a smaller function for unit testing purposes, and it should do so for
// the other methods of the Service.
func signUpService(user *User) error {
	err := user.CheckEmptyFields()
	if err != nil {
		return err
	}

	err = validateEmail(user.Email)
	if err != nil {
		return err
	}

	return nil
}

// Find a User by its ID.
func (s Service) Find(id int64) (*User, error) {
	if id == 0 {
		return &User{}, ErrUserNotFound
	}

	return s.repo.ByID(id)
}

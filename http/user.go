package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/app"
	"github.com/adrianolmedo/genesis/http/jwt"

	"github.com/gofiber/fiber/v2"
)

// loginUser handler POST: /login
func loginUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := domain.UserLoginForm{}
		err := c.BodyParser(&form)
		if err != nil {
			resp := respJSON(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.User.Login(form.Email, form.Password)
		if errors.Is(err, domain.ErrUserNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusUnauthorized).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		token, err := jwt.Generate(form.Email)
		if err != nil {
			resp := respJSON(msgError, "the token could not be generated", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		dataToken := map[string]string{"token": token}
		resp := respJSON(msgOK, "logged", dataToken)
		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// signUpUser godoc
//
//	@Summary		SignUp User
//	@Description	Register a user
//	@Accept			json
//	@Produce		json
//	@Failure		400				{object}	respError
//	@Failure		500				{object}	respError
//	@Success		201				{object}	respOkData{data=userProfileDTO}
//	@Param			userSignUpForm	body		userSignUpForm	true	"application/json"
//	@Router			/users [post]
func signUpUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := userSignUpForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(respError{
				Msg: "the JSON structure is not correct",
			})
		}

		err = s.User.SignUp(&domain.User{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(respError{
				Msg: err.Error(),
			})
		}

		return c.Status(http.StatusCreated).JSON(respOkData{
			Msg: "user created",
			Data: userProfileDTO{
				FirstName: form.FirstName,
				LastName:  form.LastName,
				Email:     form.Email,
			},
		})
	}
}

// userSignUpForm subset of User fields to create account.
type userSignUpForm struct {
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Email     string `json:"email" example:"johndoe@aol.com"`
	Password  string `json:"password" example:"1234567b"`
}

// userProfileDTO subset of User fields.
type userProfileDTO struct {
	ID        uint   `json:"id,omitempty" example:"1"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Email     string `json:"email" example:"johndoe@aol.com"`
}

// findUser handler GET: /users/:id
func findUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		user, err := s.User.Find(uint(id))
		if errors.Is(err, domain.ErrUserNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNotFound).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		resp := respJSON(msgOK, "", userProfileDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// updateUser handler PUT: /users/:id
func updateUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		form := domain.UserUpdateForm{}
		err = c.BodyParser(&form)
		if err != nil {
			resp := respJSON(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		form.ID = uint(id)

		err = s.User.Update(domain.User{
			ID:        form.ID,
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if errors.Is(err, domain.ErrUserNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		resp := respJSON(msgOK, "user updated", domain.User{
			ID:        form.ID,
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
		})

		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// listUsers handler GET: /users
func listUsers(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := s.User.List()
		if err != nil {
			resp := respJSON(msgError, "", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		if users.IsEmpty() {
			resp := respJSON(msgOK, "there are not users", nil)
			return c.Status(http.StatusOK).JSON(resp)
		}

		assemble := func(u *domain.User) userProfileDTO {
			return userProfileDTO{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
			}
		}

		list := make([]userProfileDTO, 0, len(users))
		for _, v := range users {
			list = append(list, assemble(v))
		}

		resp := respJSON(msgOK, "", list)
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// deleteUser handler DELETE: /users/:id
func deleteUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.User.Remove(int64(id))
		if errors.Is(err, domain.ErrUserNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, fmt.Sprintf("could not delete user: %s", err), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Add logger message: "User with ID %d deleted"

		resp := respJSON(msgOK, "user deleted", nil)
		return c.Status(http.StatusOK).JSON(resp)

	}
}

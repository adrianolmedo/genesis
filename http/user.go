package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	domain "github.com/adrianolmedo/genesis"
	"github.com/adrianolmedo/genesis/app"
	"github.com/adrianolmedo/genesis/http/jwt"
	"github.com/adrianolmedo/genesis/logger"

	"github.com/gofiber/fiber/v2"
)

// loginUser godoc
//
//	@Summary		Login user
//	@Description	User authentication
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Failure		400				{object}	errorResp
//	@Failure		401				{object}	errorResp
//	@Failure		500				{object}	errorResp
//	@Success		201				{object}	successResp{data=dataTokenDTO}
//	@Param			userLoginForm	body		userLoginForm	true	"application/json"
//	@Router			/login [post]
func loginUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := userLoginForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		err = s.User.Login(form.Email, form.Password)
		if errors.Is(err, domain.ErrUserNotFound) {
			return errorJSON(c, http.StatusUnauthorized, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		token, err := jwt.Generate(form.Email)
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "004",
				Message: "The token could not be generated",
			})
		}

		return successJSON(c, http.StatusCreated, respDetails{
			Message: "You are logged",
			Data:    dataTokenDTO{token},
		})
	}
}

// userLoginForm subset of user fields to request login.
type userLoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type dataTokenDTO struct {
	Token string `json:"token"`
}

// signUpUser godoc
//
//	@Summary		SignUp user
//	@Description	Register a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Failure		400				{object}	errorResp
//	@Failure		500				{object}	errorResp
//	@Success		201				{object}	successResp{data=userProfileDTO}
//	@Param			userSignUpForm	body		userSignUpForm	true	"application/json"
//	@Router			/users [post]
func signUpUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := userSignUpForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		err = s.User.SignUp(&domain.User{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		return successJSON(c, http.StatusCreated, respDetails{
			Message: "User created",
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

// findUser godoc
//
//	@Summary		Find user
//	@Description	Find user by its id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User id"
//	@Failure		400	{object}	errorResp
//	@Failure		404	{object}	errorResp
//	@Success		200	{object}	successResp{data=userProfileDTO}
//	@Router			/users/{id} [get]
func findUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}

		user, err := s.User.Find(uint(id))
		if errors.Is(err, domain.ErrUserNotFound) {
			return errorJSON(c, http.StatusNotFound, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		return successJSON(c, http.StatusOK, respDetails{
			Message: "User found",
			Data: userProfileDTO{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
		})
	}
}

// updateUser godoc
//
//	@Summary		Update user
//	@Description	Update user by its id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"User id"
//	@Failure		400				{object}	errorResp
//	@Failure		404				{object}	errorResp
//	@Success		200				{object}	successResp{data=userProfileDTO}
//	@Param			userUpdateForm	body		userUpdateForm	true	"application/json"
//	@Router			/users/{id} [put]
func updateUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}

		form := userUpdateForm{}
		err = c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		userID := uint(id)

		err = s.User.Update(domain.User{
			ID:        userID,
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if errors.Is(err, domain.ErrUserNotFound) {
			return errorJSON(c, http.StatusNotFound, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		return successJSON(c, http.StatusCreated, respDetails{
			Message: "User updated",
			Data: userProfileDTO{
				ID:        userID,
				FirstName: form.FirstName,
				LastName:  form.LastName,
				Email:     form.Email,
			},
		})
	}
}

// userUpdateForm subset of fields to request to update a User.
type userUpdateForm struct {
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Email     string `json:"email" example:"lorem@ipsum.com"`
	Password  string `json:"password" example:"1234567a"`
}

// listUsers godoc
//
//	@Summary	List users
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Failure	500	{object}	errorResp
//	@Success	200	{object}	successResp
//	@Success	200	{object}	successResp{data=[]userProfileDTO}
//	@Router		/users [get]
func listUsers(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := s.User.List()
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if users.IsEmpty() {
			return successJSON(c, http.StatusOK, respDetails{
				Code:    "005",
				Message: "There are not users",
			})
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

		return successJSON(c, http.StatusOK, respDetails{
			Message: "Ok",
			Data:    list,
		})
	}
}

// deleteUser godoc
//
//	@Summary		Delete user
//	@Description	Delete user by its id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User id"
//	@Failure		400	{object}	errorResp
//	@Failure		404	{object}	errorResp
//	@Success		200	{object}	successResp
//	@Router			/users/{id} [delete]
func deleteUser(s *app.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}

		err = s.User.Remove(uint(id))
		if errors.Is(err, domain.ErrUserNotFound) {
			return errorJSON(c, http.StatusNotFound, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: fmt.Sprintf("Could not delete user: %s", err),
			})
		}

		logger.Debug("user", fmt.Sprintf("user ID %d deleted", id))

		return successJSON(c, http.StatusOK, respDetails{
			Message: "User deleted",
		})
	}
}

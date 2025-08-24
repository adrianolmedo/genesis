package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrianolmedo/genesis/app"
	"github.com/adrianolmedo/genesis/logger"
	"github.com/adrianolmedo/genesis/pgsql"
	"github.com/adrianolmedo/genesis/rest/jwt"
	"github.com/adrianolmedo/genesis/user"

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
//	@Success		201				{object}	resp{data=dataTokenDTO}
//	@Param			userLoginForm	body		userLoginForm	true	"application/json"
//	@Router			/login [post]
func loginUser(s *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		form := userLoginForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		err = s.User.Login(ctx, form.Email, form.Password)
		if errors.Is(err, user.ErrNotFound) {
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

		return respJSON(c, http.StatusCreated, respDetails{
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
//	@Success		201				{object}	resp{data=userProfileDTO}
//	@Param			userSignUpForm	body		userSignUpForm	true	"application/json"
//	@Router			/users [post]
func signUpUser(s *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		form := userSignUpForm{}
		err := c.BodyParser(&form)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		err = s.User.SignUp(ctx, &user.User{
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

		return respJSON(c, http.StatusCreated, respDetails{
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
	ID        int64  `json:"id,omitempty" example:"1"`
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
//	@Success		200	{object}	resp{data=userProfileDTO}
//	@Router			/users/{id} [get]
func findUser(s *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}

		userModel, err := s.User.Find(ctx, int64(id))
		if errors.Is(err, user.ErrNotFound) {
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

		return respJSON(c, http.StatusOK, respDetails{
			Message: "User found",
			Data: userProfileDTO{
				ID:        userModel.ID,
				FirstName: userModel.FirstName,
				LastName:  userModel.LastName,
				Email:     userModel.Email,
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
//	@Success		200				{object}	resp{data=userProfileDTO}
//	@Param			userUpdateForm	body		userUpdateForm	true	"application/json"
//	@Router			/users/{id} [put]
func updateUser(s *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

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

		userID := int64(id)

		err = s.User.Update(ctx, user.User{
			ID:        userID,
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if errors.Is(err, user.ErrNotFound) {
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

		return respJSON(c, http.StatusCreated, respDetails{
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
//	@Success	200	{object}	resp
//	@Success	200	{object}	resp{data=[]userProfileDTO}
//	@Router		/users [get]
func listUsers(s *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		p, err := pgsql.NewPager(
			c.QueryInt("limit"),
			c.QueryInt("page"),
			c.Query("sort", "created_at"),
			c.Query("direction"),
		)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		pr, err := s.User.List(ctx, p)
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: err.Error(),
			})
		}

		users, ok := pr.Rows.(user.Users)
		if !ok {
			return errorJSON(c, http.StatusInternalServerError, respDetails{
				Code:    "003",
				Message: "Data assertion",
			})
		}

		if users.IsEmpty() {
			return respJSON(c, http.StatusOK, respDetails{
				Code:    "005",
				Message: "There are not users",
			})
		}

		assemble := func(u *user.User) userProfileDTO {
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

		return respJSON(c, http.StatusOK, respDetails{
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
//	@Success		200	{object}	resp
//	@Router			/users/{id} [delete]
func deleteUser(s *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}

		err = s.User.Remove(ctx, int64(id))
		if errors.Is(err, user.ErrNotFound) {
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

		return respJSON(c, http.StatusOK, respDetails{
			Message: "User deleted",
		})
	}
}

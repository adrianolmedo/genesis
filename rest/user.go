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
//	@Param			userLoginCommand	body	userLoginCommand	true	"application/json"
//	@Router			/login [post]
func loginUser(s *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		command := userLoginCommand{}
		err := c.BodyParser(&command)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		err = s.User.Login(ctx, command.Email, command.Password)
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

		token, err := jwt.Generate(command.Email)
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

// userLoginCommand subset of user fields to request login.
type userLoginCommand struct {
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
//	@Param			userSignUpCommand	body		userSignUpCommand	true	"application/json"
//	@Router			/users [post]
func signUpUser(s *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		command := userSignUpCommand{}
		err := c.BodyParser(&command)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, respDetails{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}

		err = s.User.SignUp(ctx, &user.User{
			FirstName: command.FirstName,
			LastName:  command.LastName,
			Email:     command.Email,
			Password:  command.Password,
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
				FirstName: command.FirstName,
				LastName:  command.LastName,
				Email:     command.Email,
			},
		})
	}
}

// userSignUpCommand subset of User fields to create account.
type userSignUpCommand struct {
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
//	@Param			userUpdateCommand	body		userUpdateCommand	true	"application/json"
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

		command := userUpdateCommand{}
		err = c.BodyParser(&command)
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
			FirstName: command.FirstName,
			LastName:  command.LastName,
			Email:     command.Email,
			Password:  command.Password,
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
				FirstName: command.FirstName,
				LastName:  command.LastName,
				Email:     command.Email,
			},
		})
	}
}

// userUpdateCommand subset of fields to request to update a User.
type userUpdateCommand struct {
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

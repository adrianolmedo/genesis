package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrianolmedo/genesis/compose"
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
//	@Success		201				{object}	resp{data=dataTokenResp}
//	@Param			userLoginReq	body		userLoginReq	true	"application/json"
//	@Router			/login [post]
func loginUser(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		req := userLoginReq{}
		err := c.BodyParser(&req)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}
		err = svcs.User.Login(ctx, req.Email, req.Password)
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("login timeout", "email", req.Email)
			return errorJSON(c, http.StatusGatewayTimeout, detailsResp{
				Code:    "005",
				Message: "The request timed out",
				Details: "Please try again later.",
			})
		}
		if errors.Is(err, user.ErrNotFound) {
			return errorJSON(c, http.StatusUnauthorized, detailsResp{
				Code:    "003",
				Message: "Invalid email or password",
			})
		}
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		token, err := jwt.Generate(req.Email)
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, detailsResp{
				Code:    "004",
				Message: "The token could not be generated",
			})
		}
		return respJSON(c, http.StatusCreated, detailsResp{
			Message: "You are logged",
			Data:    dataTokenResp{token},
		})
	}
}

// userLoginReq subset of user fields to request login.
type userLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type dataTokenResp struct {
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
//	@Success		201				{object}	resp{data=userProfileResp}
//	@Param			userSignUpReq	body		userSignUpReq	true	"application/json"
//	@Router			/users [post]
func signUpUser(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		req := userSignUpReq{}
		err := c.BodyParser(&req)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}
		err = svcs.User.SignUp(ctx, &user.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		})
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		return respJSON(c, http.StatusCreated, detailsResp{
			Message: "User created",
			Data: userProfileResp{
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Email:     req.Email,
			},
		})
	}
}

// userSignUpReq subset of User fields to create account.
type userSignUpReq struct {
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Email     string `json:"email" example:"johndoe@aol.com"`
	Password  string `json:"password" example:"1234567b"`
}

// userProfileResp subset of User fields.
type userProfileResp struct {
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
//	@Success		200	{object}	resp{data=userProfileResp}
//	@Router			/users/{id} [get]
func findUser(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}
		userModel, err := svcs.User.Find(ctx, int64(id))
		if errors.Is(err, user.ErrNotFound) {
			return errorJSON(c, http.StatusNotFound, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		return respJSON(c, http.StatusOK, detailsResp{
			Message: "User found",
			Data: userProfileResp{
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
//	@Success		200				{object}	resp{data=userProfileResp}
//	@Param			userUpdateReq	body		userUpdateReq	true	"application/json"
//	@Router			/users/{id} [put]
func updateUser(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}
		req := userUpdateReq{}
		err = c.BodyParser(&req)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "002",
				Message: "The JSON structure is not correct",
				Details: "Check the JSON syntax in the structure",
			})
		}
		userID := int64(id)
		err = svcs.User.Update(ctx, user.User{
			ID:        userID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		})
		if errors.Is(err, user.ErrNotFound) {
			return errorJSON(c, http.StatusNotFound, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		return respJSON(c, http.StatusCreated, detailsResp{
			Message: "User updated",
			Data: userProfileResp{
				ID:        userID,
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Email:     req.Email,
			},
		})
	}
}

// userUpdateReq subset of fields to request to update a User.
type userUpdateReq struct {
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
//	@Success	200	{object}	resp{data=[]userProfileResp}
//	@Router		/users [get]
func listUsers(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		p, err := pgsql.NewPager(
			c.QueryInt("limit"),
			c.QueryInt("page"),
			c.Query("sort", "created_at"),
			c.Query("direction"),
		)
		if err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		pr, err := svcs.User.List(ctx, p)
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		users, ok := pr.Rows.(user.Users)
		if !ok {
			return errorJSON(c, http.StatusInternalServerError, detailsResp{
				Code:    "003",
				Message: "Data assertion",
			})
		}
		if users.IsEmpty() {
			return respJSON(c, http.StatusOK, detailsResp{
				Code:    "005",
				Message: "There are not users",
			})
		}
		assemble := func(u user.User) userProfileResp {
			return userProfileResp{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
			}
		}
		list := make([]userProfileResp, 0, len(users))
		for _, v := range users {
			list = append(list, assemble(v))
		}
		return respJSON(c, http.StatusOK, detailsResp{
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
func deleteUser(svcs *compose.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			return errorJSON(c, http.StatusBadRequest, detailsResp{
				Code:    "002",
				Message: "Positive number expected for ID user",
			})
		}
		err = svcs.User.Remove(ctx, int64(id))
		if errors.Is(err, user.ErrNotFound) {
			return errorJSON(c, http.StatusNotFound, detailsResp{
				Code:    "003",
				Message: err.Error(),
			})
		}
		if err != nil {
			return errorJSON(c, http.StatusInternalServerError, detailsResp{
				Code:    "003",
				Message: fmt.Sprintf("Could not delete user: %s", err),
			})
		}
		logger.Debug("user", fmt.Sprintf("user ID %d deleted", id))
		return respJSON(c, http.StatusOK, detailsResp{
			Message: "User deleted",
		})
	}
}

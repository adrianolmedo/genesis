package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/service"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage"

	"github.com/labstack/echo/v4"
)

// POST: /users
func signUpUser(r storage.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		form := domain.UserSignUpForm{}

		err := c.Bind(&form)
		if err != nil {
			resp := newResponse(MsgError, "ER002", "a field in the JSON structure does not have the correct type", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = service.NewUserService(r).SignUp(domain.User{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})
		if err != nil {
			resp := newResponse(MsgError, "ER004", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(MsgOK, "OK002", "user created", domain.UserProfileDTO{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
		})

		return c.JSON(http.StatusCreated, resp)
	}
}

// GET: /users/:id
func findUser(r storage.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected for ID user", nil)
			return c.JSON(http.StatusBadRequest, resp) // 400
		}

		user, err := service.NewUserService(r).Find(int64(id))
		if errors.Is(err, domain.ErrUserNotFound) {
			resp := newResponse(MsgError, "ER007", err.Error(), nil)
			return c.JSON(http.StatusBadRequest, resp) // http.StatusNoContent 204
		}

		if err != nil {
			resp := newResponse(MsgError, "ER009", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		//u.Password = ""
		resp := newResponse(MsgOK, "OK002", "", domain.UserProfileDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
		return c.JSON(http.StatusOK, resp)
	}
}

// PUT: /users/:id.
func updateUser(r storage.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected for ID user", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		form := domain.UserUpdateForm{}

		err = c.Bind(&form)
		if err != nil {
			resp := newResponse(MsgError, "ER002", "a field in the JSON structure does not have the correct type", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		form.ID = int64(id)

		err = service.NewUserService(r).Update(domain.User{
			ID:        form.ID,
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})
		if errors.Is(err, domain.ErrUserNotFound) {
			resp := newResponse(MsgError, "ER007", err.Error(), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(MsgError, "ER004", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(MsgOK, "OK002", "user updated", domain.UserProfileDTO{
			ID:        form.ID,
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
		})
		return c.JSON(http.StatusOK, resp)
	}
}

// listUsers handler for GET: /users.
func listUsers(r storage.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := service.NewUserService(r).List()
		if err != nil {
			resp := newResponse(MsgError, "ER003", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		if len(users) == 0 {
			resp := newResponse(MsgOK, "OK002", "there are not users", nil)
			return c.JSON(http.StatusOK, resp) // maybe 204
		}

		list := make(domain.UsersList, 0, len(users))

		assemble := func(u *domain.User) domain.UserProfileDTO {
			return domain.UserProfileDTO{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
			}
		}

		for _, v := range users {
			list = append(list, assemble(v))
		}

		resp := newResponse(MsgOK, "OK002", "", list)
		return c.JSON(http.StatusOK, resp)
	}
}

// deleteUser handler for DELETE: /users/:id.
func deleteUser(r storage.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected for ID user", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = service.NewUserService(r).Remove(int64(id))
		if errors.Is(err, domain.ErrUserNotFound) {
			resp := newResponse(MsgError, "ER007", err.Error(), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(MsgError, "ER006", fmt.Sprintf("could not delete user: %s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(MsgOK, "OK002", "user deleted", nil)
		return c.JSON(http.StatusOK, resp) // maybe 204
	}
}

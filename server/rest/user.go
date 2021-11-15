package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go-restapi-practice/user"

	"github.com/labstack/echo/v4"
)

// createUser handler for POST: /users.
func createUser(r user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := user.ProfileRequest{}

		err := c.Bind(&req)
		if err != nil {
			resp := newResponse(MsgError, "ER002", "a field in the JSON structure does not have the correct type", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = user.NewService(r).Register(&user.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		})
		if err != nil {
			resp := newResponse(MsgError, "ER004", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(MsgOK, "OK002", "resource created", user.ProfileResponse{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		})

		return c.JSON(http.StatusCreated, resp)
	}
}

// getUserByID handler for GET: /users/:id.
func getUserByID(r user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected", nil)
			return c.JSON(http.StatusBadRequest, resp) // 400
		}

		data, err := user.NewService(r).ByID(int64(id))
		if errors.Is(err, user.ErrResourceDoesNotExist) {
			resp := newResponse(MsgError, "ER007", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(MsgError, "ER007", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		//data.Password = ""
		resp := newResponse(MsgOK, "OK002", "", user.ProfileResponse{
			ID:        data.ID,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
		})
		return c.JSON(http.StatusOK, resp)
	}
}

// updateUserByID handler for PUT: /users/:id.
func updateUserByID(r user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		req := user.UpdateRequest{}
		err = c.Bind(&req)
		if err != nil {
			resp := newResponse(MsgError, "ER002", "a field in the JSON structure does not have the correct type", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		req.ID = int64(id)

		err = user.NewService(r).Update(&user.User{
			ID:        req.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		})
		if errors.Is(err, user.ErrResourceDoesNotExist) {
			resp := newResponse(MsgError, "ER007", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(MsgError, "ER004", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(MsgOK, "OK002", "resource updated", user.ProfileResponse{
			ID:        req.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		})
		return c.JSON(http.StatusOK, resp)
	}
}

// getAllUsers handler for GET: /users.
func getAllUsers(r user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := user.NewService(r).All()
		if err != nil {
			resp := newResponse(MsgError, "ER003", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		if len(users) == 0 {
			resp := newResponse(MsgOK, "OK002", "there are not resources", nil)
			return c.JSON(http.StatusOK, resp) // maybe 204
		}

		list := make(user.ListResponse, 0, len(users))

		assemble := func(m *user.User) user.ProfileResponse {
			return user.ProfileResponse{
				ID:        m.ID,
				FirstName: m.FirstName,
				LastName:  m.LastName,
				Email:     m.Email,
			}
		}

		for _, v := range users {
			list = append(list, assemble(v))
		}

		resp := newResponse(MsgOK, "OK002", "", list)
		return c.JSON(http.StatusOK, resp)
	}
}

// deleteUserByID handler for DELETE: /users/:id.
func deleteUserByID(r user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = user.NewService(r).Delete(int64(id))
		if errors.Is(err, user.ErrResourceDoesNotExist) {
			resp := newResponse(MsgError, "ER007", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(MsgError, "ER006", fmt.Sprintf("could not delete resource: %s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(MsgOK, "OK002", "resource deleted", nil)
		return c.JSON(http.StatusOK, resp) // maybe 204
	}
}

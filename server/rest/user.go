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
func createUser(u user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := user.Request{}

		err := c.Bind(&data)
		if err != nil {
			resp := newResponse(MsgError, "ER002", "a field in the JSON structure does not have the correct type", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = user.NewService(u).Register(&data)
		if err != nil {
			resp := newResponse(MsgError, "ER004", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		data.Password = ""
		resp := newResponse(MsgOK, "OK002", "resource created", data)
		return c.JSON(http.StatusCreated, resp)
	}
}

// getUserByID handler for GET: /users/:id.
func getUserByID(u user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected", nil)
			return c.JSON(http.StatusBadRequest, resp) // 400
		}

		data, err := user.NewService(u).ByID(int64(id))
		if errors.Is(err, user.ErrResourceDoesNotExist) {
			resp := newResponse(MsgError, "ER007", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(MsgError, "ER007", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		//data.Password = ""
		resp := newResponse(MsgOK, "OK002", "", data)
		return c.JSON(http.StatusOK, resp)
	}
}

// updateUserByID handler for PUT: /users/:id.
func updateUserByID(u user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		data := user.Request{}
		err = c.Bind(&data)
		if err != nil {
			resp := newResponse(MsgError, "ER002", "a field in the JSON structure does not have the correct type", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		data.ID = int64(id)
		err = user.NewService(u).Update(&data)
		if errors.Is(err, user.ErrResourceDoesNotExist) {
			resp := newResponse(MsgError, "ER007", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusNoContent, resp)
		}

		if err != nil {
			resp := newResponse(MsgError, "ER004", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp := newResponse(MsgOK, "OK002", "resource updated", data)
		return c.JSON(http.StatusOK, resp)
	}
}

// getAllUsers handler for GET: /users.
func getAllUsers(u user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := user.NewService(u).All()
		if err != nil {
			resp := newResponse(MsgError, "ER003", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		if len(users) == 0 {
			resp := newResponse(MsgOK, "OK002", "there are not resources", nil)
			return c.JSON(http.StatusOK, resp) // maybe 204
		}

		resp := newResponse(MsgOK, "OK002", "", users)
		return c.JSON(http.StatusOK, resp)
	}
}

// deleteUserByID handler for DELETE: /users/:id.
func deleteUserByID(u user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if id <= 0 || err != nil {
			resp := newResponse(MsgError, "ER005", "positive number expected", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = user.NewService(u).Delete(int64(id))
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

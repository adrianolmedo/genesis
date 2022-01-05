package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/adrianolmedo/go-restapi-practice/infra/jwt"
	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/service"
	"github.com/adrianolmedo/go-restapi-practice/internal/storage"

	"github.com/labstack/echo/v4"
)

// POST: /login
func loginUser(r storage.LoginRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		form := domain.UserLoginForm{}

		// Bind nos ayuda hacer el volcado a una estructura JSON,
		// automáticamente Bind captura el r.Body o w del ResponseWriter.
		err := c.Bind(&form)
		if err != nil {
			resp := newResponse(MsgError, "ER002", "a field in the JSON structure does not have the correct type", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = service.NewLoginService(r).Execute(form.Email, form.Password)
		if errors.Is(err, domain.ErrUserNotFound) {
			resp := newResponse(MsgError, "ER007", err.Error(), nil)
			return c.JSON(http.StatusUnauthorized, resp)
		}

		if err != nil {
			resp := newResponse(MsgError, "ER009", fmt.Sprintf("%s", err), nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		token, err := jwt.New(form.Email)
		if err != nil {
			resp := newResponse(MsgError, "ER008", "the token could not be generated", nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		dataToken := map[string]string{"token": token}
		resp := newResponse(MsgOK, "OK004", "logged", dataToken)
		return c.JSON(http.StatusCreated, resp)
	}
}

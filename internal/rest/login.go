package rest

import (
	"errors"
	"net/http"

	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
	"github.com/adrianolmedo/go-restapi-practice/internal/service"
	"github.com/adrianolmedo/go-restapi-practice/jwt"

	"github.com/labstack/echo/v4"
)

// loginUser handler POST: /login
func loginUser(s service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		form := domain.UserLoginForm{}

		// Bind nos ayuda hacer el volcado a una estructura JSON,
		// automáticamente Bind captura el r.Body o w del ResponseWriter.
		err := c.Bind(&form)
		if err != nil {
			resp := newResponse(msgError, "ER002", "the JSON structure is not correct", nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = s.Login.Execute(form.Email, form.Password)
		if errors.Is(err, domain.ErrUserNotFound) {
			resp := newResponse(msgError, "ER007", err.Error(), nil)
			return c.JSON(http.StatusUnauthorized, resp)
		}

		if err != nil {
			resp := newResponse(msgError, "ER009", err.Error(), nil)
			return c.JSON(http.StatusBadRequest, resp)
		}

		token, err := jwt.Generate(form.Email)
		if err != nil {
			resp := newResponse(msgError, "ER008", "the token could not be generated", nil)
			return c.JSON(http.StatusInternalServerError, resp)
		}

		dataToken := map[string]string{"token": token}
		resp := newResponse(msgOK, "OK004", "logged", dataToken)
		return c.JSON(http.StatusCreated, resp)
	}
}

package rest

import (
	"net/http"

	"go-restapi-practice/auth"
	"go-restapi-practice/model"

	"github.com/labstack/echo/v4"
)

// login for validate JWT. POST: /login.
func login(c echo.Context) error {
	data := model.Login{}

	// Bind nos ayuda hacer el volcado a una estructura JSON,
	// automáticamente Bind captura el r.Body o w del ResponseWriter.
	err := c.Bind(&data)
	if err != nil {
		resp := newResponse(MsgError, "ER002", "a field in the JSON structure does not have the correct type", nil)
		return c.JSON(http.StatusBadRequest, resp)
	}

	if !validateLogin(&data) {
		resp := newResponse(MsgError, "ER007", "invalid username and password", nil)
		return c.JSON(http.StatusBadRequest, resp)
	}

	token, err := auth.GenerateToken(&data)
	if err != nil {
		resp := newResponse(MsgError, "ER008", "the token could not be generated", nil)
		return c.JSON(http.StatusInternalServerError, resp)
	}

	dataToken := map[string]string{"token": token}
	resp := newResponse(MsgOK, "OK004", "logged", dataToken)
	return c.JSON(http.StatusCreated, resp)
}

// validateLogin is a bussisnes logic for validation username and password.
func validateLogin(data *model.Login) bool {
	return data.Email == "adrianolmedo@gmail.com" && data.Password == "1234567@"
}

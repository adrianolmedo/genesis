package rest

import "github.com/labstack/echo/v4"

// reqMeth simulate func signing of HTTP request methods of Echo framework
// like GET, POST, PUT, DELETE, etc.
func reqMeth(c echo.Context, h echo.HandlerFunc) error {
	return h(c)
}

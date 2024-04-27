package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Auth = middleware.BasicAuth(checkAuth)

func checkAuth(u, p string, c echo.Context) (bool, error) {
	au := os.Getenv("ADMIN_USERNAME")
	ap := os.Getenv("ADMIN_PASSWORD")

	if au == "" || ap == "" {
		return false, echo.ErrInternalServerError
	}
	return u == au && p == ap, nil
}

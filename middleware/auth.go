package middleware

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Auth = middleware.BasicAuth(checkAuth)

func checkAuth(u, p string, c echo.Context) (bool, error) {
	fmt.Println("🚀 | file: auth.go | line 14 | funccheckAuth | p : ", p)
	fmt.Println("🚀 | file: auth.go | line 14 | funccheckAuth | u : ", u)
	au := os.Getenv("ADMIN_USERNAME")
	fmt.Println("🚀 | file: auth.go | line 14 | funccheckAuth | au : ", au)
	ap := os.Getenv("ADMIN_PASSWORD")
	fmt.Println("🚀 | file: auth.go | line 16 | funccheckAuth | ap : ", ap)

	if au == "" || ap == "" {
		return false, echo.ErrInternalServerError
	}
	return u == au && p == ap, nil
}

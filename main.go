package main

import (
	"net/http"

	"github.com/jaiieth/assessment-tax/calculator"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/middleware"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger)
	e.Validator = helper.NewValidator()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.POST("/tax/calculations", calculator.Handler)

	e.Logger.Fatal(e.Start(":1323"))
}

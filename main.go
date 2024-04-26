package main

import (
	"net/http"

	"github.com/jaiieth/assessment-tax/calculator"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/middleware"
	"github.com/jaiieth/assessment-tax/postgres"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		panic("failed to connect database")
	}
	e := echo.New()

	e.Use(middleware.Logger)
	e.Validator = helper.NewValidator()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	handler := calculator.New(db)

	e.POST("/tax/calculations", handler.CalculateTax)

	e.Logger.Fatal(e.Start(":1323"))
}

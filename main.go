package main

import (
	"net/http"

	"github.com/jaiieth/assessment-tax/config"
	"github.com/jaiieth/assessment-tax/handler"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := config.New()
	if err != nil {
		panic("failed to connect database")
	}
	e := echo.New()

	e.Use(middleware.Logger)
	e.Validator = helper.NewValidator()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	h := handler.New(db)

	e.GET("/tax/config", h.GetConfig)
	e.POST("/tax/calculations", h.CalculateTax)
	e.POST("/tax/calculations/upload-csv", h.CalculateByCsv)
	e.POST("/admin/deductions/personal", h.SetPersonalDeduction)

	e.Logger.Fatal(e.Start(":1323"))
}

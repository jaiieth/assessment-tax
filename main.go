package main

import (
	"net/http"

	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/middleware"
	"github.com/jaiieth/assessment-tax/pkg/calculator"
	"github.com/jaiieth/assessment-tax/pkg/config"
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

	c := calculator.NewHandler(db)
	e.POST("/tax/calculations", c.CalculateTaxHandler)
	e.POST("/tax/calculations/upload-csv", c.CalculateByCsvHandler)

	a := config.NewHandler(db)
	admin := e.Group("/admin", middleware.Auth)
	admin.GET("/config", a.GetConfigHandler)
	admin.POST("/deductions/personal", a.SetPersonalDeductionHandler)
	admin.POST("/deductions/k-receipt", a.SetMaxKReceiptHandler)

	e.Logger.Fatal(e.Start(":1323"))
}

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
	admin := e.Group("/admin", middleware.Auth)

	e.GET("/tax/config", h.GetConfigHandler)
	e.POST("/tax/calculations", h.CalculateTaxHandler)
	e.POST("/tax/calculations/upload-csv", h.CalculateByCsvHandler)

	admin.POST("/deductions/personal", h.SetPersonalDeductionHandler)

	e.Logger.Fatal(e.Start(":1323"))
}

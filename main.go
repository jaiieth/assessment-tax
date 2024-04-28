package main

import (
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

	admin := e.Group("/admin", middleware.Auth)

	c := calculator.NewHandler(db)
	a := config.NewHandler(db)

	c.RegisterRoutes(e)
	a.RegisterRoutes(admin)

	e.Logger.Fatal(e.Start(":1323"))
}

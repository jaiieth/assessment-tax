package main

import (
	"fmt"
	"os"

	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/middleware"
	"github.com/jaiieth/assessment-tax/pkg/calculator"
	"github.com/jaiieth/assessment-tax/pkg/config"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	//Load env
	godotenv.Load()

	//Init DB
	db, err := config.New()
	if err != nil {
		panic("failed to connect database")
	}

	//Init Echo
	e := echo.New()
	port := os.Getenv("PORT")

	e.Use(middleware.Logger)
	e.Validator = helper.NewValidator()

	admin := e.Group("/admin", middleware.Auth)

	c := calculator.NewHandler(db)
	a := config.NewHandler(db)

	c.RegisterRoutes(e)
	a.RegisterRoutes(admin)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}

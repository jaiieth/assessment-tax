package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(fmt.Sprintf(":%v", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(fmt.Sprintf("err: failed to shutdown server %v", err))
	}
	log.Println("Server stopped")
}

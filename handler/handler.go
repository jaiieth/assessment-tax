package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jaiieth/assessment-tax/handler/calculator"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/postgres"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB Database
}

type Database interface {
	GetConfig() (calculator.Config, error)
	SetPersonalDeduction(float64) (calculator.Config, error)
}

var validate *validator.Validate

func New(db Database) Handler {
	return Handler{DB: db}
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func (h Handler) CalculateTax(c echo.Context) error {
	var body calculator.CalculateTaxBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: "invalid request"})
	}

	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: "invalid request"})
	}

	config, err := h.DB.GetConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{Message: "Oops, something went wrong"})
	}

	res := calculator.CalculateTax(body, config)

	return c.JSON(http.StatusOK, res)

}

func (h Handler) SetPersonalDeduction(c echo.Context) error {

	var body calculator.PersonalDeductionBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: "invalid request"})
	}

	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: "invalid request"})
	}

	if body.Amount > postgres.MAX_PERSONAL_DEDUCTION {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: fmt.Sprintf("invalid request: Maximum personal deduction is %0.f", postgres.MAX_PERSONAL_DEDUCTION)})
	}
	if body.Amount < 10000 {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: fmt.Sprintf("invalid request: Minimum personal deduction is %0.f", postgres.MIN_PERSONAL_DEDUCTION)})
	}

	config, err := h.DB.SetPersonalDeduction(body.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, config)
	}
	return c.JSON(http.StatusOK, config)
}

func (h Handler) GetConfig(c echo.Context) error {
	config, err := h.DB.GetConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{Message: "Oops, something went wrong"})

	}
	return c.JSON(http.StatusOK, config)
}

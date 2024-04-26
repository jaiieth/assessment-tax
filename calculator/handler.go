package calculator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB Database
}

type Database interface {
	GetConfig() (Config, error)
}

var validate *validator.Validate

func New(db Database) Handler {
	return Handler{DB: db}
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func (h Handler) CalculateTax(c echo.Context) error {
	var body CalculateTaxBody
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

	res := CalculateTax(body, config)

	return c.JSON(http.StatusOK, res)

}

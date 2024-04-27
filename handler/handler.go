package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jaiieth/assessment-tax/config"
	calc "github.com/jaiieth/assessment-tax/handler/calculator"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB Database
}

type Database interface {
	GetConfig() (config.Config, error)
	SetPersonalDeduction(float64) (config.Config, error)
}

var validate *validator.Validate

func New(db Database) Handler {
	return Handler{DB: db}
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func (h Handler) CalculateTaxHandler(c echo.Context) error {
	var body calc.CalculateTaxBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	config, err := h.DB.GetConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorRes("Oops, something went wrong"))
	}

	res := calc.CalculateTax(body, config)

	return c.JSON(http.StatusOK, res)

}

func (h Handler) SetPersonalDeductionHandler(c echo.Context) error {
	var body calc.SetPersonalDeductionBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	if body.Amount > config.MAX_PERSONAL_DEDUCTION {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes(fmt.Sprintf("invalid request: Maximum personal deduction is %0.f", config.MAX_PERSONAL_DEDUCTION)))
	}
	if body.Amount < 10000 {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes(fmt.Sprintf("invalid request: Maximum personal deduction is %0.f", config.MIN_PERSONAL_DEDUCTION)))
	}

	config, err := h.DB.SetPersonalDeduction(body.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, config)
	}
	return c.JSON(http.StatusOK, config)
}

func (h Handler) GetConfigHandler(c echo.Context) error {
	config, err := h.DB.GetConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorRes("Oops, something went wrong"))

	}
	return c.JSON(http.StatusOK, config)
}

func (h Handler) CalculateByCsvHandler(c echo.Context) error {
	file, err := c.FormFile("taxes.csv")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	src, err := file.Open()
	if err != nil {
		fmt.Println("🚀 | file: handler.go | line 105 | func | err : ", err)
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}
	defer src.Close()

	i := TaxCSVInstance{src}

	err = i.validate()
	if err != nil {
		fmt.Println("🚀 | file: handler.go | line 113 | func | err : ", err)
		return c.JSON(http.StatusBadRequest, helper.ErrorRes(err.Error()))
	}

	var records []calc.TaxCSV
	if err := i.unmarshal(&records); err != nil {
		fmt.Println("🚀 | file: handler.go | line 119 | iferr:=i.unmarshal | err : ", err)
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	for _, r := range records {
		if err := validate.Struct(r); err != nil {
			fmt.Println("🚀 | file: handler.go | line 125 | func | err : ", err)
			return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
		}
	}

	config, err := h.DB.GetConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorRes("Oops, something went wrong"))
	}

	res := calc.CalculateTaxByCsv(records, config)
	return c.JSON(http.StatusOK, calc.CalculateByCSVResponse{Taxes: res})
}

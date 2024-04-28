package calculator

import (
	"fmt"
	"net/http"

	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/pkg/config"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB config.Database
}

func NewHandler(db config.Database) Handler {
	return Handler{DB: db}
}

func (h Handler) CalculateTaxHandler(c echo.Context) error {
	var body CalculateTaxBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	if err := c.Validate(body); err != nil {
		fmt.Println("🚀 | file: handler.go | line 34 | iferr:=c.Validate | err : ", err)
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	config, err := h.DB.GetConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorRes("Oops, something went wrong"))
	}

	res := CalculateTax(body, config)

	return c.JSON(http.StatusOK, res)

}

func (h Handler) CalculateByCsvHandler(c echo.Context) error {
	file, err := c.FormFile("taxes.csv")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}
	defer src.Close()

	i := TaxCSVInstance{src}

	err = i.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes(err.Error()))
	}

	var records []TaxCSV
	if err := i.Unmarshal(&records); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	for _, r := range records {
		if err := c.Validate(r); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
		}
	}

	config, err := h.DB.GetConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorRes("Oops, something went wrong"))
	}

	res := CalculateTaxes(records, config)
	return c.JSON(http.StatusOK, CalculateByCSVResponse{Taxes: res})
}

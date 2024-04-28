package config

import (
	"fmt"
	"net/http"

	"github.com/jaiieth/assessment-tax/helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB Database
}

func NewHandler(db Database) Handler {
	return Handler{DB: db}
}

func (h Handler) SetPersonalDeductionHandler(c echo.Context) error {
	var d Deduction

	if err := d.BindAndValidateStruct(c); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	if err := d.ValidateValue(MIN_PERSONAL_DEDUCTION, MAX_PERSONAL_DEDUCTION); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes(fmt.Sprintf(
			"Personal deduction must be between %0.f and %0.f",
			MAX_PERSONAL_DEDUCTION, MIN_PERSONAL_DEDUCTION,
		)))
	}

	config, err := h.DB.SetPersonalDeduction(*d.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, config)
	}
	return c.JSON(http.StatusOK, config)
}

func (h Handler) SetMaxKReceiptHandler(c echo.Context) error {
	var d Deduction

	if err := d.BindAndValidateStruct(c); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes("invalid request"))
	}

	if err := d.ValidateValue(MIN_K_RECEIPT, MAX_K_RECEIPT); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorRes(fmt.Sprintf(
			"Maximum K-Receipt must be between  %0.f and %0.f",
			MAX_K_RECEIPT, MIN_K_RECEIPT,
		)))
	}

	config, err := h.DB.SetMaxKReceipt(*d.Amount)
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

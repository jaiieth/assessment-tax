package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
	"github.com/jaiieth/assessment-tax/handler/calculator"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/postgres/config"
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

func (h Handler) CalculateTax(c echo.Context) error {
	var body calculator.CalculateTaxBody
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

	res := calculator.CalculateTax(body, config)

	return c.JSON(http.StatusOK, res)

}

func (h Handler) SetPersonalDeduction(c echo.Context) error {
	var body calculator.SetPersonalDeductionBody
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

func (h Handler) GetConfig(c echo.Context) error {
	config, err := h.DB.GetConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorRes("Oops, something went wrong"))

	}
	return c.JSON(http.StatusOK, config)
}

type TaxCSV struct {
	TotalIncome    float64  `csv:"totalIncome" validate:"required,numeric,gte=0"`
	WithHoldingTax *float64 `csv:"wht" validate:"gte=0"`
	Donation       *float64 `csv:"donation" validate:"gte=0"`
}

func (h Handler) CalculateByCsv(c echo.Context) error {
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

	var records []TaxCSV
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

	return c.JSON(http.StatusOK, records)
}

type TaxCSVInstance struct {
	File multipart.File
}

func (ti TaxCSVInstance) validate() error {
	rows, err := gocsv.LazyCSVReader(ti.File).ReadAll()
	if err != nil {
		return err
	}

	header := rows[0]
	expectedHeaders := []string{"totalIncome", "wht", "donation"}

	for i, h := range header {
		if h != expectedHeaders[i] {
			return fmt.Errorf("wrong csv format")
		}
	}

	for _, row := range rows {
		for _, value := range row {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("wrong csv format")
			}
		}
	}
	// Rewind to the beginning of csv, So the `t.File` can be read again
	ti.File.Seek(0, 0)
	return nil
}

func (t TaxCSVInstance) unmarshal(s interface{}) error {
	if err := gocsv.UnmarshalMultipartFile(&t.File, s); err != nil {
		return err
	}

	return nil
}

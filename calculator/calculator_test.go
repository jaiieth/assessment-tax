package calculator_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jaiieth/assessment-tax/calculator"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type GetTotalTaxCases struct {
	name        string
	input       float64
	expectedTax float64
}

func NewContext(method string, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()

	e.Validator = helper.NewValidator()

	context := e.NewContext(req, rec)
	return context, rec
}

func RunTestGetTotalTax(t *testing.T, cases []GetTotalTaxCases) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tax := calculator.GetTotalTax(c.input)
			fmt.Println("🚀 | file: calculator_test.go | line 70 | t.Run | tax : ", tax)
			assert.Equal(t, c.expectedTax, tax)
		})
	}
}

func TestGetTotalTax(t *testing.T) {
	cases := []GetTotalTaxCases{
		{name: "Given income 150000 should return 0", input: 150000, expectedTax: 0},
		{name: "Given income 150,001 should return 0.1", input: 150001, expectedTax: 0.1},
		{name: "Given income 500,000 should return 35,000", input: 500000, expectedTax: 35000},
		{name: "Given income 500,001 should return 35,000.15", input: 500001, expectedTax: 35000.15},
		{name: "Given income 1,000,000 should return 35,000 + 75,0000 = 110,000", input: 1000000, expectedTax: 110000},
		{name: "Given income 1,000,001 should return 35,000 + 75,0000 = 110,000.2", input: 1000001, expectedTax: 110000.2},
		{name: "Given income 2,000,000 should return 110,000 + 200,000 = 310,000", input: 2000000, expectedTax: 310000},
		{name: "Given income 2,000,001 should return 110,000 + 200,000 = 310,000.3", input: 2000001, expectedTax: 310000.35},
		{name: "Given income 3,000,000 should return 310,000 + 350,000 = 660,000", input: 3000000, expectedTax: 660000},
	}

	RunTestGetTotalTax(t, cases)
}

func TestCalculateTax(t *testing.T) {
	t.Run("Given income 0 with WHT should return tax:0 and taxRefund:WHT", func(t *testing.T) {
		body := calculator.CalculateTaxBody{
			TotalIncome:    0,
			WithHoldingTax: 50000,
		}

		tax, taxRefund := calculator.CalculateTax(body)
		assert.Equal(t, float64(0), tax)
		assert.Equal(t, 50000.0, taxRefund)
	})
	t.Run("Given income 500,000 with no WHT should return tax:29,000 and taxRefund:0 ", func(t *testing.T) {
		body := calculator.CalculateTaxBody{
			TotalIncome:    500000,
			WithHoldingTax: 0.0,
		}

		tax, _ := calculator.CalculateTax(body)
		assert.Equal(t, 29000.0, tax)
	})

	t.Run("Given income 500,000 with 25,000 WHT should return tax:4000 and taxRefund:0", func(t *testing.T) {
		body := calculator.CalculateTaxBody{
			TotalIncome:    500000,
			WithHoldingTax: 25000,
		}

		tax, _ := calculator.CalculateTax(body)
		assert.Equal(t, 4000.0, tax)
	})
}

func TestCalculationHandler(t *testing.T) {
	t.Run("Given valid request body should return 200", func(t *testing.T) {

		c, rec := NewContext(http.MethodPost, "/tax/calculations", strings.NewReader(`
		{
			"totalIncome": 500000.0,
			"wht": 0.0,
			"allowances": [
				{
					"allowanceType": "donation",
					"amount": 0.0
				}
			]
		}`))
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		err := calculator.Handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response calculator.CalculateResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	})

	t.Run("Given invalid request body should return 400", func(t *testing.T) {
		c, rec := NewContext(http.MethodPost, "/tax/calculations", strings.NewReader(`{"totalIncome": Invalid}`))
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		err := calculator.Handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response calculator.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	})

	t.Run("Given no request body should return 400", func(t *testing.T) {
		c, rec := NewContext(http.MethodPost, "/tax/calculations", strings.NewReader(`{}`))
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		err := calculator.Handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response calculator.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "invalid request", response.Message)
	})
}

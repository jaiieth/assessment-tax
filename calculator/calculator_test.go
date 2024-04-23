package calculator_test

import (
	"bytes"
	"encoding/json"
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
	input       float32
	expectedTax float32
}

func NewContext(method string, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()

	e.Validator = helper.NewValidator()

	context := e.NewContext(req, rec)
	return context, rec
}

func RunTestCalculateTaxWithNoAllowance(t *testing.T, tests []GetTotalTaxCases) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := calculator.CalculateTaxBody{
				TotalIncome:    test.input,
				WithHoldingTax: 0.0,
			}

			body, err := json.Marshal(input)

			if err != nil {
				t.Error(err)
			}
			c, rec := NewContext(http.MethodPost, "/tax/calculations", bytes.NewBuffer(body))
			c.Request().Header.Set("Content-Type", "application/json")

			handler := calculator.New()

			handler.CalculationHandler(c)

			got := calculator.CalculateResponse{}
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, test.expectedTax, got.Tax)
		})

	}
}

func RunTestGetTotalTax(t *testing.T, cases []GetTotalTaxCases) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tax := calculator.GetTotalTax(c.input)
			assert.Equal(t, c.expectedTax, tax)
		})
	}
}

func TestGetTotalTax(t *testing.T) {
	cases := []GetTotalTaxCases{
		{name: "Given income 150000 should return 0", input: 150000, expectedTax: 0},
		{name: "Given income 150,001 should return 0.1", input: 150001, expectedTax: 0.1},
		{name: "Given income 500,000 should return 35,000", input: 500000, expectedTax: 35000},
		{name: "Given income 1,000,000 should return 35,000+ 75,0000", input: 1000000, expectedTax: 110000},
		{name: "Given income 2,000,000 should return 110,000 + 200,000", input: 2000000, expectedTax: 310000},
		{name: "Given income 3,000,000 should return 310,000 + 350,000", input: 3000000, expectedTax: 660000},
	}

	RunTestGetTotalTax(t, cases)
}

func TestCalculate(t *testing.T) {
	h := calculator.New()
	res := h.CalculateTax(210000)
	assert.Equal(t, float32(0.0), res)
	res = h.CalculateTax(210001)
	assert.Greater(t, res, float32(0.0))
}

func TestCalculationHandler(t *testing.T) {
	t.Run("Given valid request body should return 200", func(t *testing.T) {

		c, rec := NewContext(http.MethodPost, "/tax/calculations", strings.NewReader(`{"totalIncome": 50000}`))
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		h := &calculator.Handler{}

		err := h.CalculationHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response calculator.CalculateResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	})

	t.Run("Given invalid request body should return 400", func(t *testing.T) {
		c, rec := NewContext(http.MethodPost, "/tax/calculations", strings.NewReader(`{"totalIncome": Invalid}`))
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		h := &calculator.Handler{}

		err := h.CalculationHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response calculator.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	})

	t.Run("Given no request body should return 400", func(t *testing.T) {
		c, rec := NewContext(http.MethodPost, "/tax/calculations", strings.NewReader(`{}`))
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		h := &calculator.Handler{}

		err := h.CalculationHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response calculator.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "invalid request", response.Message)
	})
}

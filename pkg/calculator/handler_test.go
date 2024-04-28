package calculator_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaiieth/assessment-tax/helper"
	calc "github.com/jaiieth/assessment-tax/pkg/calculator"
	"github.com/jaiieth/assessment-tax/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	Config config.Config
	Error  error
	mock.Mock
}

func (m *mockDB) GetConfig() (config.Config, error) {
	return m.Config, m.Error
}
func (m *mockDB) SetPersonalDeduction(n float64) (config.Config, error) {
	m.Called(n)

	return m.Config, nil
}
func (m *mockDB) SetMaxKReceipt(n float64) (config.Config, error) {
	m.Called()
	return m.Config, nil
}

func TestCalculateTaxHandler(t *testing.T) {
	t.Run("TestSuccessfulRequestWithValidInput", func(t *testing.T) {
		body := calc.CalculateTaxBody{
			TotalIncome:    500000,
			WithHoldingTax: 50000,
			Allowances:     []calc.Allowance{{Type: "donation", Amount: 50000}},
		}
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Errorf("failed to marshal body: %v", err)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		e.Validator = helper.NewValidator()
		c := e.NewContext(req, rec)

		h := calc.NewHandler(&mockDB{})
		h.CalculateTaxHandler(c)

		assert.Equal(t, http.StatusOK, rec.Code, fmt.Sprintf("status code should be %d but got %v", http.StatusOK, rec.Code))

		var response calc.CalculateTaxResult
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "failed to unmarshal response body")
	})

	t.Run("TestInvalidRequestWithMissingInputValues", func(t *testing.T) {
		body := calc.CalculateTaxBody{
			WithHoldingTax: 50000,
			Allowances:     []calc.Allowance{{Type: "donation", Amount: 50000}},
		}
		bodyJSON, _ := json.Marshal(body)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		e.Validator = helper.NewValidator()
		c := e.NewContext(req, rec)

		h := calc.NewHandler(&mockDB{})
		h.CalculateTaxHandler(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code, fmt.Sprintf("status code should be %d but got %v", http.StatusOK, rec.Code))
	})

	t.Run("TestInvalidRequestWithInvalidInput", func(t *testing.T) {
		body := calc.CalculateTaxBody{
			TotalIncome:    -5000,
			WithHoldingTax: 50000,
			Allowances:     []calc.Allowance{{Type: "donation", Amount: 50000}},
		}
		bodyJSON, _ := json.Marshal(body)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		e.Validator = helper.NewValidator()
		c := e.NewContext(req, rec)

		h := calc.NewHandler(&mockDB{})
		h.CalculateTaxHandler(c)

		var response helper.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rec.Code, fmt.Sprintf("status code should be %d but got %v", http.StatusOK, rec.Code))
		assert.NoError(t, err, "failed to unmarshal response body")
		assert.Equal(t, helper.ErrorRes("invalid request"), response)
	})

	t.Run("TestInvalidRequestWithInvalidInput_InvalidJSONBody", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer([]byte(`{Invalid}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		e.Validator = helper.NewValidator()
		c := e.NewContext(req, rec)

		h := calc.NewHandler(&mockDB{})
		h.CalculateTaxHandler(c)

		var response helper.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rec.Code, fmt.Sprintf("status code should be %d but got %v", http.StatusOK, rec.Code))
		assert.NoError(t, err, "failed to unmarshal response body")
		assert.Equal(t, helper.ErrorRes("invalid request"), response)
	})
}

// Successful request with valid input

func TestErrorWhenUnableToRetrieveConfig(t *testing.T) {
	body := calc.CalculateTaxBody{
		TotalIncome: 500000,
	}
	bodyJSON, _ := json.Marshal(body)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{Error: errors.New("failed to retrieve config")}
	h := calc.NewHandler(db)
	h.DB = db

	h.CalculateTaxHandler(c)

	var response helper.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NoError(t, err, "failed to unmarshal response body")
	assert.Equal(t, helper.ErrorRes("Oops, something went wrong"), response)
}

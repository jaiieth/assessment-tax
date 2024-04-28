package calculator_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

// Valid CSV file is uploaded and processed successfully
func TestValidCSVFile(t *testing.T) {

	// Create a temporary file to mimic a real file upload
	file, err := os.CreateTemp("", "taxes.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name()) // Ensure the temporary file is deleted
	defer file.Close()           // Ensure the file is closed after writing

	csvData := []byte("totalIncome,wht,donation\n10000,500,200\n20000,1000,400\n")
	if _, err := file.Write(csvData); err != nil {
		t.Fatal(err)
	}

	// Rewind to the start of the file so it can be read during the request handling
	_, err = file.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare a multipart writer
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile("taxes.csv", file.Name())
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err)
	}

	mw.Close()

	// Reset file position for reading
	file.Seek(0, 0)

	e := echo.New()
	e.Validator = helper.NewValidator()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", &b)
	req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
	c := e.NewContext(req, rec)

	h := calc.NewHandler(&mockDB{})
	h.CalculateByCsvHandler(c)

	var res calc.CalculateByCSVResponse

	assert.Equal(t, http.StatusOK, rec.Code, "expected status code %d, got %d", http.StatusOK, rec.Code)
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res), "failed to unmarshal response body")
}

// Valid CSV file is uploaded and processed successfully
func TestGetConfigError(t *testing.T) {

	// Create a temporary file to mimic a real file upload
	file, err := os.CreateTemp("", "taxes.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name()) // Ensure the temporary file is deleted
	defer file.Close()           // Ensure the file is closed after writing

	csvData := []byte("totalIncome,wht,donation\n10000,500,200\n20000,1000,400\n")
	if _, err := file.Write(csvData); err != nil {
		t.Fatal(err)
	}

	// Rewind to the start of the file so it can be read during the request handling
	_, err = file.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare a multipart writer
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile("taxes.csv", file.Name())
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err)
	}

	mw.Close()

	// Reset file position for reading
	file.Seek(0, 0)

	e := echo.New()
	e.Validator = helper.NewValidator()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", &b)
	req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
	c := e.NewContext(req, rec)

	stubDB := &mockDB{
		Error: errors.New("failed to retrieve config"),
	}

	h := calc.NewHandler(stubDB)
	h.CalculateByCsvHandler(c)

	var res calc.CalculateByCSVResponse

	assert.Equal(t, http.StatusInternalServerError, rec.Code, "expected status code %d, got %d", http.StatusOK, rec.Code)
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res), "failed to unmarshal response body")
}
func TestNoTaxesCSVFile(t *testing.T) {
	// Create a temporary file to mimic a real file upload
	file, err := os.CreateTemp("", "taxes.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name()) // Ensure the temporary file is deleted
	defer file.Close()           // Ensure the file is closed after writing
	// Reset file position for reading
	file.Seek(0, 0)

	e := echo.New()
	e.Validator = helper.NewValidator()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	c := e.NewContext(req, rec)

	h := calc.NewHandler(&mockDB{})
	h.CalculateByCsvHandler(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d, got %d", http.StatusOK, rec.Code)
}

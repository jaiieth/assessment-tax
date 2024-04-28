package calculator_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	csvData := []byte("TotalIncome,Donation,WithHoldingTax\n10000,500,200\n20000,1000,400\n")
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
	_, err = mw.CreateFormFile("file", "taxes.csv")
	if err != nil {
		t.Fatal(err)
	}

	mw.Close()

	// Reset file position for reading
	file.Seek(0, 0)

	e := echo.New()
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

// // Empty CSV file is uploaded
// func TestEmptyCSVFile(t *testing.T) {
// 	// Mocking the echo.Context
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Creating a temporary empty CSV file
// 	file, err := ioutil.TempFile("", "taxes.csv")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(file.Name())

// 	// Setting up the form file in the request context
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
// 	req.MultipartForm = &multipart.Form{
// 		File: map[string][]*multipart.FileHeader{
// 			"taxes.csv": []*multipart.FileHeader{
// 				{
// 					Filename: "taxes.csv",
// 					Size:     0,
// 				},
// 			},
// 		},
// 	}

// 	// Calling the handler function
// 	if err := h.CalculateByCsvHandler(c); err != nil {
// 		t.Fatal(err)
// 	}

// 	// Asserting the response status code
// 	if rec.Code != http.StatusBadRequest {
// 		t.Errorf("expected status code %d but got %d", http.StatusBadRequest, rec.Code)
// 	}

// 	// Asserting the response body
// 	var res ErrorResponse
// 	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
// 		t.Fatal(err)
// 	}

// 	expectedRes := ErrorResponse{
// 		Message: "invalid request",
// 	}

// 	if !reflect.DeepEqual(res, expectedRes) {
// 		t.Errorf("expected response %+v but got %+v", expectedRes, res)
// 	}
// }

// // Personal deduction and donation are within valid range
// func TestValidCSVFile(t *testing.T) {
// 	// Mocking the echo.Context
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Creating a temporary CSV file
// 	file, err := ioutil.TempFile("", "taxes.csv")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(file.Name())

// 	// Writing valid CSV data to the file
// 	csvData := []byte("TotalIncome,Donation,WithHoldingTax\n10000,500,200\n20000,1000,400\n")
// 	if _, err := file.Write(csvData); err != nil {
// 		t.Fatal(err)
// 	}

// 	// Setting up the form file in the request context
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
// 	req.MultipartForm = &multipart.Form{
// 		File: map[string][]*multipart.FileHeader{
// 			"taxes.csv": []*multipart.FileHeader{
// 				{
// 					Filename: "taxes.csv",
// 					Size:     int64(len(csvData)),
// 				},
// 			},
// 		},
// 	}

// 	// Calling the handler function
// 	if err := h.CalculateByCsvHandler(c); err != nil {
// 		t.Fatal(err)
// 	}

// 	// Asserting the response status code
// 	if rec.Code != http.StatusOK {
// 		t.Errorf("expected status code %d but got %d", http.StatusOK, rec.Code)
// 	}

// 	// Asserting the response body
// 	var res CalculateByCSVResponse
// 	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
// 		t.Fatal(err)
// 	}

// 	expectedRes := CalculateByCSVResponse{
// 		Taxes: []CalculateByCSVResponseItem{
// 			{TotalIncome: 10000, Tax: 2300, Refund: 0},
// 			{TotalIncome: 20000, Tax: 4600, Refund: 0},
// 		},
// 	}

// 	if !reflect.DeepEqual(res, expectedRes) {
// 		t.Errorf("expected response %+v but got %+v", expectedRes, res)
// 	}
// }

// // Personal deduction or donation exceed valid range
// func TestInvalidPersonalDeductionOrDonation(t *testing.T) {
// 	// Mocking the echo.Context
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Creating a temporary CSV file
// 	file, err := ioutil.TempFile("", "taxes.csv")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(file.Name())

// 	// Writing invalid CSV data to the file
// 	csvData := []byte("TotalIncome,Donation,WithHoldingTax\n10000,1500,200\n20000,3000,400\n")
// 	if _, err := file.Write(csvData); err != nil {
// 		t.Fatal(err)
// 	}

// 	// Setting up the form file in the request context
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
// 	req.MultipartForm = &multipart.Form{
// 		File: map[string][]*multipart.FileHeader{
// 			"taxes.csv": []*multipart.FileHeader{
// 				{
// 					Filename: "taxes.csv",
// 					Size:     int64(len(csvData)),
// 				},
// 			},
// 		},
// 	}

// 	// Calling the handler function
// 	if err := h.CalculateByCsvHandler(c); err != nil {
// 		t.Fatal(err)
// 	}

// 	// Asserting the response status code
// 	if rec.Code != http.StatusBadRequest {
// 		t.Errorf("expected status code %d but got %d", http.StatusBadRequest, rec.Code)
// 	}

// 	// Asserting the response body
// 	var res ErrorResponse
// 	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
// 		t.Fatal(err)
// 	}

// 	expectedRes := ErrorResponse{
// 		Message: "invalid request",
// 	}

// 	if !reflect.DeepEqual(res, expectedRes) {
// 		t.Errorf("expected response %+v but got %+v", expectedRes, res)
// 	}
// }

// // Negative withholding tax is provided
// func TestNegativeWithholdingTax(t *testing.T) {
// 	// Mocking the echo.Context
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Creating a temporary CSV file
// 	file, err := ioutil.TempFile("", "taxes.csv")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(file.Name())

// 	// Writing invalid CSV data to the file
// 	csvData := []byte("TotalIncome,Donation,WithHoldingTax\n10000,500,-200\n20000,1000,-400\n")
// 	if _, err := file.Write(csvData); err != nil {
// 		t.Fatal(err)
// 	}

// 	// Setting up the form file in the request context
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
// 	req.MultipartForm = &multipart.Form{
// 		File: map[string][]*multipart.FileHeader{
// 			"taxes.csv": []*multipart.FileHeader{
// 				{
// 					Filename: "taxes.csv",
// 					Size:     int64(len(csvData)),
// 				},
// 			},
// 		},
// 	}

// 	handler:= calc.NewHandler(&mockDB{})
// 	// Calling the handler function
// 	if err := h.CalculateByCsvHandler(c); err != nil {
// 		t.Fatal(err)
// 	}

// 	// Asserting the response status code
// 	if rec.Code != http.StatusBadRequest {
// 		t.Errorf("expected status code %d but got %d", http.StatusBadRequest, rec.Code)
// 	}

// 	// Asserting the response body
// 	var res helper.ErrorResponse
// 	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
// 		t.Fatal(err)
// 	}

// 	expectedRes := helper.ErrorResponse{
// 		Message: "invalid request",
// 	}

// 	if !reflect.DeepEqual(res, expectedRes) {
// 		t.Errorf("expected response %+v but got %+v", expectedRes, res)
// 	}
// }

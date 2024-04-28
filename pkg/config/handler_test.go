package config_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaiieth/assessment-tax/helper"
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
	return m.Config, m.Error
}
func (m *mockDB) SetMaxKReceipt(n float64) (config.Config, error) {
	m.Called(n)
	return m.Config, m.Error
}
func TestSetPersonalDeductionHandler_ValidInput(t *testing.T) {
	body := config.Deduction{
		Amount: float64Ptr(50000.0),
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal body: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{
		Config: config.Config{
			PersonalDeduction: 50000.0,
		},
	}

	h := config.NewHandler(db)
	db.On("SetPersonalDeduction", 50000.0).Return()
	h.SetPersonalDeductionHandler(c)
	db.AssertCalled(t, "SetPersonalDeduction", 50000.0)

	err = json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		t.Errorf("response is not JSON")
	}
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 50000.0, *body.Amount)
}

func TestSetPersonalDeductionHandler_InvalidInput(t *testing.T) {
	body := config.Deduction{}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal body: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{
		Config: config.Config{
			PersonalDeduction: 50000.0,
		},
	}

	h := config.NewHandler(db)
	h.SetPersonalDeductionHandler(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSetPersonalDeductionHandler_ValueNotInLimit(t *testing.T) {
	body := config.Deduction{
		Amount: float64Ptr(150000.0),
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal body: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{
		Config: config.Config{
			PersonalDeduction: 50000.0,
		},
	}

	h := config.NewHandler(db)
	h.SetPersonalDeductionHandler(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSetPersonalDeductionHandler_GetConfigError(t *testing.T) {
	body := config.Deduction{
		Amount: float64Ptr(50000.0),
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal body: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{
		Error: errors.New("failed to get config"),
	}

	h := config.NewHandler(db)
	db.On("SetPersonalDeduction", 50000.0).Return()
	h.SetPersonalDeductionHandler(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestSetMaxKReceiptHandler_ValidInput(t *testing.T) {
	body := config.Deduction{
		Amount: float64Ptr(50000.0),
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal body: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{
		Config: config.Config{
			PersonalDeduction: 50000.0,
		},
	}

	h := config.NewHandler(db)
	db.On("SetMaxKReceipt", 50000.0).Return()
	h.SetMaxKReceiptHandler(c)
	db.AssertCalled(t, "SetMaxKReceipt", 50000.0)

	err = json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		t.Errorf("response is not JSON")
	}
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 50000.0, *body.Amount)
}

func TestSetMaxKReceiptHandler_InvalidInput(t *testing.T) {
	body := config.Deduction{}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal body: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{
		Config: config.Config{
			PersonalDeduction: 50000.0,
		},
	}

	h := config.NewHandler(db)
	h.SetMaxKReceiptHandler(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSetMaxKReceiptHandler_ValueNotInLimit(t *testing.T) {
	body := config.Deduction{
		Amount: float64Ptr(150000.0),
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal body: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{
		Config: config.Config{
			PersonalDeduction: 50000.0,
		},
	}

	h := config.NewHandler(db)
	h.SetMaxKReceiptHandler(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSetMaxKReceiptHandler_GetConfigError(t *testing.T) {
	body := config.Deduction{
		Amount: float64Ptr(50000.0),
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal body: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	e.Validator = helper.NewValidator()
	c := e.NewContext(req, rec)

	db := &mockDB{
		Error: errors.New("failed to get config"),
	}

	h := config.NewHandler(db)
	db.On("SetMaxKReceipt", 50000.0).Return()
	h.SetMaxKReceiptHandler(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetConfigHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		e := echo.New()
		c := e.NewContext(req, rec)

		db := &mockDB{
			Config: config.Config{
				PersonalDeduction: 40000.0,
				MaxKReceipt:       50000.0,
			}}

		h := config.NewHandler(db)

		h.GetConfigHandler(c)
		assert.Equal(t, http.StatusOK, rec.Code)

		var body config.Config
		err := json.Unmarshal(rec.Body.Bytes(), &body)
		if err != nil {
			t.Errorf("response is not JSON")
		}
		assert.Equal(t, 40000.0, body.PersonalDeduction)
		assert.Equal(t, 50000.0, body.MaxKReceipt)
	})
	t.Run("Failed", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		e := echo.New()
		c := e.NewContext(req, rec)

		db := &mockDB{
			Error: errors.New("failed to get config"),
		}

		h := config.NewHandler(db)

		h.GetConfigHandler(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}

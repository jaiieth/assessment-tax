package config_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaiieth/assessment-tax/helper"
	"github.com/jaiieth/assessment-tax/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	t.Run("Success", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows([]string{"personal_deduction", "max_k_receipt"}).AddRow(5000, 10000))

		p := &config.Postgres{
			Db: db,
		}

		expPersonalDeduction := 5000.0
		expMaxKReceipt := 10000.0

		config, err := p.GetConfig()

		assert.NoError(t, err)
		assert.Equal(t, expPersonalDeduction, config.PersonalDeduction)
		assert.Equal(t, expMaxKReceipt, config.MaxKReceipt)
	})

	t.Run("Failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT *").WillReturnError(sql.ErrNoRows)

		p := &config.Postgres{
			Db: db,
		}

		_, err = p.GetConfig()

		assert.Error(t, err)
	})
}

func TestSetMaxKRecepit(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		queryArgs := 10000.0
		mock.ExpectQuery("UPDATE config").
			WithArgs(queryArgs).WillReturnRows(sqlmock.NewRows([]string{"max_k_receipt"}).AddRow(queryArgs))

		p := &config.Postgres{
			Db: db,
		}

		expectedResult := 10000.0

		config, err := p.SetMaxKReceipt(queryArgs)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, config.MaxKReceipt)
	})

	t.Run("Failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("UPDATE config").WillReturnError(sql.ErrNoRows)

		p := &config.Postgres{
			Db: db,
		}

		_, err = p.SetMaxKReceipt(10000.0)

		assert.Error(t, err)
	})
}

func TestSetPersonalDeduction(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		queryArgs := 10000.0
		mock.ExpectQuery("UPDATE config").
			WithArgs(queryArgs).WillReturnRows(sqlmock.NewRows([]string{"personal_deduction"}).AddRow(queryArgs))

		p := &config.Postgres{
			Db: db,
		}

		expectedResult := 10000.0

		config, err := p.SetPersonalDeduction(queryArgs)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, config.PersonalDeduction)
	})

	t.Run("Failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("UPDATE config").WillReturnError(sql.ErrNoRows)

		p := &config.Postgres{
			Db: db,
		}

		_, err = p.SetPersonalDeduction(10000.0)

		assert.Error(t, err)
	})
}

func TestBindAndValidateStruct(t *testing.T) {
	// The function should correctly bind and validate the JSON request body.
	t.Run("ValidRequestBody", func(t *testing.T) {
		e := echo.New()
		e.Validator = helper.NewValidator()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": 0}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		d := &config.Deduction{}

		err := d.BindAndValidateStruct(c)

		assert.NoError(t, err)
	})

	// The function should return an error when the JSON request body is valid but value less than 0.
	t.Run("ValidButValueLessThanZero", func(t *testing.T) {
		e := echo.New()
		e.Validator = helper.NewValidator()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": -1}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		d := config.Deduction{}

		err := d.BindAndValidateStruct(c)

		assert.Error(t, err)
	})

	// The function should return an error when the JSON request body is empty.
	t.Run("EmptyRequestBody", func(t *testing.T) {
		e := echo.New()
		e.Validator = helper.NewValidator()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		d := config.Deduction{}

		err := d.BindAndValidateStruct(c)

		assert.Error(t, err)
	})

	// The function should return an error when the JSON request body is invalid.
	t.Run("InvalidRequestBody", func(t *testing.T) {
		e := echo.New()
		e.Validator = helper.NewValidator()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": "invalid"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		d := config.Deduction{}

		err := d.BindAndValidateStruct(c)

		assert.Error(t, err)
	})
}

func TestValidateValue(t *testing.T) {
	t.Run("Given amount within range should return nil", func(t *testing.T) {
		deduction := config.Deduction{
			Amount: float64Ptr(50000.0),
		}

		err := deduction.ValidateValue(0.0, 100000.0)
		assert.NoError(t, err)
	})

	t.Run("Given amount equal maximum limit should return nil", func(t *testing.T) {
		deduction := config.Deduction{
			Amount: float64Ptr(100000.0),
		}

		err := deduction.ValidateValue(0.0, 100000.0)
		assert.NoError(t, err)
	})

	t.Run("Given amount equal minimum limit should return nil", func(t *testing.T) {
		deduction := config.Deduction{
			Amount: float64Ptr(0.0),
		}

		err := deduction.ValidateValue(0.0, 100000.0)
		assert.NoError(t, err)
	})

	t.Run("Given amount more than maximum limit should return error", func(t *testing.T) {
		deduction := config.Deduction{
			Amount: float64Ptr(100001.0),
		}

		err := deduction.ValidateValue(0.0, 100000.0)
		assert.Error(t, err)
	})

	t.Run("Given amount less than minimum limit should return error", func(t *testing.T) {
		deduction := config.Deduction{
			Amount: float64Ptr(9999.0),
		}

		err := deduction.ValidateValue(10000.0, 100000.0)
		assert.Error(t, err)
	})

}

func float64Ptr(f float64) *float64 {
	return &f
}

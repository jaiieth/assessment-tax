package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCheckAuth(t *testing.T) {
	t.Run("Invalid ENV", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		data := "admin:password"
		encoded := base64.StdEncoding.EncodeToString([]byte(data))

		req.Header.Set("Authorization", "Basic "+encoded)
		c := e.NewContext(req, rec)

		isValid, err := checkAuth("admin", "password", c)

		assert.Error(t, err)
		assert.False(t, isValid)
	})

	os.Setenv("ADMIN_USERNAME", "admin")
	os.Setenv("ADMIN_PASSWORD", "password")

	t.Run("Valid credentials", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		data := "admin:password"
		encoded := base64.StdEncoding.EncodeToString([]byte(data))

		req.Header.Set("Authorization", "Basic "+encoded)
		c := e.NewContext(req, rec)

		isValid, err := checkAuth("admin", "password", c)

		assert.NoError(t, err)
		assert.True(t, isValid)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		username := os.Getenv("ADMIN_USERNAME")
		password := os.Getenv("ADMIN_PASSWORD")
		data := fmt.Sprintf("%s:%s", username, password)
		encoded := base64.StdEncoding.EncodeToString([]byte(data))

		req.Header.Set("Authorization", "Basic "+encoded)
		c := e.NewContext(req, rec)

		isValid, err := checkAuth(username, "wrong password", c)

		assert.NoError(t, err)
		assert.False(t, isValid)
	})
}

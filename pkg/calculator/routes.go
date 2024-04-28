package calculator

import "github.com/labstack/echo/v4"

func (h Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/tax/calculations", h.CalculateTaxHandler)
	e.POST("/tax/calculations/upload-csv", h.CalculateByCsvHandler)
}

package config

import "github.com/labstack/echo/v4"

func (h Handler) RegisterRoutes(e *echo.Group) {
	e.GET("/config", h.GetConfigHandler)
	e.POST("/deductions/personal", h.SetPersonalDeductionHandler)
	e.POST("/deductions/k-receipt", h.SetMaxKReceiptHandler)
}

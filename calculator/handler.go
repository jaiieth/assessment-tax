package calculator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	PersonalDeduction float64
}

func New() *Handler {
	return &Handler{PersonalDeduction: PersonalDeduction}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
type SuccessResponse struct {
	Message string `json:"message"`
}

type CalculateResponse struct {
	Tax float64 `json:"tax"`
}

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func (h *Handler) CalculateTax(totalIncome float64) float64 {
	return GetTotalTax(totalIncome - h.PersonalDeduction)

}
func (h *Handler) CalculationHandler(c echo.Context) error {
	var body CalculateTaxBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
	}

	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
	}

	tax := h.CalculateTax(body.TotalIncome)

	return c.JSON(http.StatusOK, CalculateResponse{Tax: tax})

}

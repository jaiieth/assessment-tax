package calculator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
}
type SuccessResponse struct {
	Message string `json:"message"`
}

type CalculateResponse struct {
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"taxRefund,omitempty"`
}

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func Handler(c echo.Context) error {
	var body CalculateTaxBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
	}

	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
	}

	res := CalculateTax(body)

	return c.JSON(http.StatusOK, res)

}

package helper

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

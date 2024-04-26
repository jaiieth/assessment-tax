package calculator

type Config struct {
	PersonalDeduction float64 `json:"personalDeduction,omitempty"  validate:"gte=0"`
	MaxDonation       float64 `json:"maxDonation,omitempty"`
}

type PersonalDeductionBody struct {
	Amount float64 `json:"amount" validate:"gte=0"`
}

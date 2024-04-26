package calculator

type Config struct {
	PersonalDeduction float64 `json:"personalDeduction" validate:"gte=0"`
	MaxDonation       float64
}

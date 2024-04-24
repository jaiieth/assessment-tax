package calculator

import "math"

type CalculateTaxBody struct {
	TotalIncome    float64     `json:"totalIncome" validate:"required,gte=0"`
	WithHoldingTax float64     `json:"wht" validate:"gte=0"`
	Allowances     []Allowance `json:"allowances"`
}

type Allowance struct {
	Type   string  `json:"allowanceType"  example:"donation" validate:"required ,oneof='donation' 'k-receipt'"`
	Amount float64 `json:"amount" validate:"required,gte=0"`
}

const (
	Donation = "donation"
	KReceipt = "k-receipt"
)

var AllowanceType = []string{
	Donation, KReceipt,
}

var PersonalDeduction float64 = 60000.0

func GetTotalTax(taxable float64) float64 {
	if taxable > 2000000 {
		return ((taxable - 2000000) * 0.35) + GetTotalTax(2000000)
	}

	if taxable > 1000000 {
		return ((taxable - 1000000) * 0.20) + GetTotalTax(1000000)
	}

	if taxable > 500000 {
		return ((taxable - 500000) * 0.15) + GetTotalTax(500000)
	}

	if taxable > 150000 {
		return ((taxable - 150000) * 0.10) + GetTotalTax(150000)
	}

	return 0.0
}

func CalculateTax(b CalculateTaxBody) (tax, taxRefund float64) {
	tax = GetTotalTax(b.TotalIncome-PersonalDeduction) - b.WithHoldingTax

	if tax < 0 {
		return 0, math.Abs(tax)
	}
	return tax, 0
}

package calculator

import (
	"math"

	"github.com/jaiieth/assessment-tax/config"
)

type CalculateTaxBody struct {
	TotalIncome    float64     `json:"totalIncome" validate:"required,gte=0"`
	WithHoldingTax float64     `json:"wht" validate:"gte=0"`
	Allowances     []Allowance `json:"allowances" validate:"unique=Type,dive"`
}

type SetPersonalDeductionBody struct {
	Amount float64 `json:"amount" validate:"gte=0"`
}

type Allowance struct {
	Type   string  `json:"allowanceType"  example:"donation" validate:"required,oneof=donation k-receipt"`
	Amount float64 `json:"amount" validate:"gte=0"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type CalculateTaxResponse struct {
	Tax       float64    `json:"tax"`
	TaxLevel  []TaxLevel `json:"taxLevel,omitempty"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
}

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

func GetTaxLevels(taxable float64) (taxLevel []TaxLevel) {
	taxLevel = append(taxLevel, TaxLevel{Level: "0-150,000", Tax: 0})
	taxable = math.Max(taxable-150000, 0)
	taxLevel = append(taxLevel, TaxLevel{Level: "150,001-500,000", Tax: (math.Min(taxable, 350000) * 0.10)})
	taxable = math.Max(taxable-350000, 0)
	taxLevel = append(taxLevel, TaxLevel{Level: "500,001-1,000,000", Tax: (math.Min(taxable, 500000) * 0.15)})
	taxable = math.Max(taxable-500000, 0)
	taxLevel = append(taxLevel, TaxLevel{Level: "1,000,001-2,000,000", Tax: (math.Min(taxable, 1000000) * 0.20)})
	taxable = math.Max(taxable-1000000, 0)
	taxLevel = append(taxLevel, TaxLevel{Level: "2,000,001 ขึ้นไป", Tax: (math.Min(taxable, 2000000) * 0.35)})

	return taxLevel
}

func getDonationAllowance(allowances []Allowance) float64 {
	donation := 0.0
	maxDonation := config.MAX_DONATION

	for _, a := range allowances {
		if a.Type == config.AllowanceType.Donation {
			donation += a.Amount
		}
	}

	return math.Min(donation, maxDonation)
}

func CalculateTax(b CalculateTaxBody, c config.Config) CalculateTaxResponse {
	allowance := getDonationAllowance(b.Allowances)

	tax := GetTotalTax(b.TotalIncome-c.PersonalDeduction-allowance) - b.WithHoldingTax
	roundedTax := math.Round(tax*100) / 100
	var taxLevel []TaxLevel
	if tax < 0 {
		return CalculateTaxResponse{0, taxLevel, math.Abs(roundedTax)}
	}

	taxLevel = GetTaxLevels(b.TotalIncome - allowance - c.PersonalDeduction)
	return CalculateTaxResponse{math.Max(0, roundedTax), taxLevel, 0}
}

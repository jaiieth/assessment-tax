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

type CalculateTaxResult struct {
	Tax       float64    `json:"tax"`
	TaxLevel  []TaxLevel `json:"taxLevel,omitempty"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
}

type TaxCSV struct {
	TotalIncome    float64  `csv:"totalIncome" validate:"required,numeric,gte=0"`
	WithHoldingTax *float64 `csv:"wht" validate:"gte=0"`      //use pointer to allow 0
	Donation       *float64 `csv:"donation" validate:"gte=0"` // use pointer to allow 0
}

type CalculateByCSVResponse struct {
	Taxes []CalculateByCSVResponseItem `json:"taxes"`
}

type CalculateByCSVResponseItem struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
	TaxRefund   float64 `json:"taxRefund,omitempty"`
}

func roundTwoDigits(x float64) float64 {
	return math.Round(x*100) / 100
}
func GetTotalTax(taxable float64) float64 {

	if taxable > 2000000 {
		return roundTwoDigits((taxable-2000000)*0.35) + GetTotalTax(2000000)
	}

	if taxable > 1000000 {
		return roundTwoDigits((taxable-1000000)*0.20) + GetTotalTax(1000000)
	}

	if taxable > 500000 {
		return roundTwoDigits((taxable-500000)*0.15) + GetTotalTax(500000)
	}

	if taxable > 150000 {
		return roundTwoDigits((taxable-150000)*0.10) + GetTotalTax(150000)
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

func getDonationAllowance(allowances []Allowance) (donation float64) {
	for _, a := range allowances {
		if a.Type == config.AllowanceType.Donation {
			donation += a.Amount
		}
	}

	return math.Min(donation, config.MAX_DONATION)
}

func CalculateTax(b CalculateTaxBody, c config.Config) CalculateTaxResult {
	allowance := getDonationAllowance(b.Allowances)

	tax := GetTotalTax(b.TotalIncome-c.PersonalDeduction-allowance) - b.WithHoldingTax
	var taxLevel []TaxLevel
	if tax < 0 {
		return CalculateTaxResult{0, taxLevel, math.Abs(tax)}
	}

	taxLevel = GetTaxLevels(b.TotalIncome - allowance - c.PersonalDeduction)
	return CalculateTaxResult{math.Max(0, tax), taxLevel, 0}
}

func CalculateTaxByCsv(rs []TaxCSV, c config.Config) []CalculateByCSVResponseItem {
	res := []CalculateByCSVResponseItem{}
	for _, r := range rs {
		allowance := math.Min(*r.Donation, config.MAX_DONATION)
		tax := GetTotalTax(r.TotalIncome-c.PersonalDeduction-allowance) - *r.WithHoldingTax
		if tax < 0 {
			res = append(res, CalculateByCSVResponseItem{r.TotalIncome, 0, math.Abs(tax)})
			continue
		}

		res = append(res, CalculateByCSVResponseItem{r.TotalIncome, tax, 0})
	}

	return res
}

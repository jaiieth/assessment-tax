package config

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Config struct {
	PersonalDeduction float64 `postgres:"personal_deduction" json:"personalDeduction,omitempty"`
	MaxKReceipt       float64 `postgres:"max_k_receipt" json:"kReceipt,omitempty"`
}

type Database interface {
	GetConfig() (Config, error)
	SetPersonalDeduction(float64) (Config, error)
	SetMaxKReceipt(float64) (Config, error)
}

const (
	DEFAULT_PERSONAL_DEDUCTION = 60000.0
	DEFAULT_MAX_K_RECEIPT      = 50000.0
	MAX_K_RECEIPT              = 100000.0
	MIN_K_RECEIPT              = 0.0
	MAX_DONATION               = 100000.0
	MAX_PERSONAL_DEDUCTION     = 100000.0
	MIN_PERSONAL_DEDUCTION     = 10000.0
)

var AllowanceType = struct {
	Donation string
	KReceipt string
}{
	Donation: "donation",
	KReceipt: "k-receipt",
}

func (p *Postgres) GetConfig() (c Config, err error) {
	c = Config{
		PersonalDeduction: DEFAULT_PERSONAL_DEDUCTION,
		MaxKReceipt:       DEFAULT_MAX_K_RECEIPT,
	}
	err = p.Db.QueryRow("SELECT * FROM config").Scan(&c.PersonalDeduction, &c.MaxKReceipt)

	if err != nil {
		return Config{}, err
	}

	return c, nil

}

func (p *Postgres) SetPersonalDeduction(n float64) (config Config, err error) {
	err = p.Db.QueryRow("UPDATE config SET personal_deduction = $1 RETURNING personal_deduction", n).Scan(&config.PersonalDeduction)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (p *Postgres) SetMaxKReceipt(n float64) (config Config, err error) {
	err = p.Db.QueryRow("UPDATE config SET max_k_receipt = $1 RETURNING max_k_receipt", n).Scan(&config.MaxKReceipt)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

type Deduction struct {
	Amount *float64 `json:"amount" validate:"required,gte=0"`
}

func (d Deduction) BindAndValidateStruct(c echo.Context) error {
	if err := c.Bind(&d); err != nil {
		return fmt.Errorf("err: Cannot bind JSON")
	}

	if err := c.Validate(d); err != nil {
		return fmt.Errorf("err: Invailid json body")
	}

	return nil
}

func (d Deduction) ValidateValue(min float64, max float64) error {
	if *d.Amount < min || *d.Amount > max {
		return fmt.Errorf("err: Not in range")
	}
	return nil
}

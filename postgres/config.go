package postgres

import (
	"fmt"

	"github.com/jaiieth/assessment-tax/handler/calculator"
)

const (
	DEFAULT_PERSONAL_DEDUCTION = 60000.0
	MAX_PERSONAL_DEDUCTION     = 100000.0
	MIN_PERSONAL_DEDUCTION     = 10000.0
	MAX_DONATION               = 100000.0
)

type Config struct {
	PersonalDeduction float64 `postgres:"personal_deduction" json:"personalDeduction"`
	MaxDonation       float64 `postgres:"max_donation" json:"maxDonation"`
}

func (p *Postgres) GetConfig() (config calculator.Config, err error) {
	config = calculator.Config{
		PersonalDeduction: DEFAULT_PERSONAL_DEDUCTION,
		MaxDonation:       MAX_DONATION,
	}
	err = p.Db.QueryRow("SELECT * FROM config").Scan(&config.PersonalDeduction)

	if err != nil {
		fmt.Println("🚀 | file: config.go | line 25 | func | err : ", err)
		return calculator.Config{}, err
	}

	return config, nil

}

func (p *Postgres) SetPersonalDeduction(n float64) (config calculator.Config, err error) {
	err = p.Db.QueryRow("UPDATE config SET personal_deduction = $1 RETURNING personal_deduction", n).Scan(&config.PersonalDeduction)
	if err != nil {
		return calculator.Config{}, err
	}
	return config, nil
}

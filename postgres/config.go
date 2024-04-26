package postgres

import "github.com/jaiieth/assessment-tax/calculator"

const (
	DEFAULT_PERSONAL_DEDUCTION = 60000.0
	MAX_PERSONAL_DEDUCTION     = 100000.0
	MIN_PERSONAL_DEDUCTION     = 10000.0
	MAX_DONATION               = 100000.0
)

type Config struct {
	PersonalDeduction float64 `postgres:"personal_deduction"`
	MaxDonation       float64 `postgres:"max_donation"`
}

func (p *Postgres) GetConfig() (config calculator.Config, err error) {
	config = calculator.Config{
		PersonalDeduction: DEFAULT_PERSONAL_DEDUCTION,
		MaxDonation:       MAX_DONATION,
	}
	err = p.Db.QueryRow("SELECT * FROM config").Scan(&config.PersonalDeduction)

	if err != nil {
		return calculator.Config{}, err
	}

	return config, nil

}

func (p *Postgres) MaximumDonation() (float64, error) {
	return 0.0, nil
}

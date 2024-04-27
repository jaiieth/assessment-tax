package config

import (
	"database/sql"
)

type Config struct {
	PersonalDeduction float64 `postgres:"personal_deduction" json:"personalDeduction"`
	MaxDonation       float64 `postgres:"max_donation" json:"maxDonation"`
}

type Postgres struct {
	Db *sql.DB
}

func (p *Postgres) GetConfig() (c Config, err error) {
	c = Config{
		PersonalDeduction: DEFAULT_PERSONAL_DEDUCTION,
		MaxDonation:       MAX_DONATION,
	}
	err = p.Db.QueryRow("SELECT * FROM config").Scan(&c.PersonalDeduction)

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

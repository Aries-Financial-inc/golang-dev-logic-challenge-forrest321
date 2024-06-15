package models

import (
	"errors"
	"strings"
	"time"
)

const (
	Call  = "call"
	Put   = "put"
	Long  = "long"
	Short = "short"
)

// OptionsContract holds the details of a given Contract for analysis
type OptionsContract struct {
	Type           string    `json:"type"`
	StrikePrice    float64   `json:"strike_price"`
	Bid            float64   `json:"bid"`
	Ask            float64   `json:"ask"`
	ExpirationDate time.Time `json:"expiration_date"`
	Position       string    `json:"long_short"`
}

// GetType ensures that the text is in an expected case
func (o OptionsContract) GetType() string {
	return strings.ToLower(o.Type)
}

// GetPosition ensures that the text is in an expected case
func (o OptionsContract) GetPosition() string {
	return strings.ToLower(o.Position)
}

// Validate ensures that the contract has valid properties
func (o OptionsContract) Validate() []error {
	var errs []error

	if strings.ToLower(o.Type) != Call && strings.ToLower(o.Type) != Put {
		errs = []error{errors.New("invalid type")}
	}
	if o.StrikePrice <= 0 {
		errs = append(errs, errors.New("invalid strike price"))
	}
	if o.Bid <= 0 {
		errs = append(errs, errors.New("invalid bid"))
	}
	if o.Ask <= 0 {
		errs = append(errs, errors.New("invalid ask"))
	}
	if o.ExpirationDate.IsZero() || o.ExpirationDate.Before(time.Now()) {
		errs = append(errs, errors.New("invalid expiration date"))
	}
	if strings.ToLower(o.Position) != Long && strings.ToLower(o.Position) != Short {
		errs = append(errs, errors.New("invalid long short"))
	}

	return errs
}

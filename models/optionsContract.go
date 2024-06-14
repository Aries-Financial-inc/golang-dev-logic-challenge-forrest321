package models

import "time"

type OptionType string

const (
	Call OptionType = "call"
	Put  OptionType = "put"
)

type PositionType string

const (
	Long  PositionType = "long"
	Short PositionType = "short"
)

type OptionsContract struct {
	Type           OptionType   `json:"type"`
	StrikePrice    float64      `json:"strike_price"`
	Bid            float64      `json:"bid"`
	Ask            float64      `json:"ask"`
	ExpirationDate time.Time    `json:"expiration_date"`
	Position       PositionType `json:"long_short"`
}

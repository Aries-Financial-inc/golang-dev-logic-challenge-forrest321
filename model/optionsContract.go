package model

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
	Type           OptionType
	StrikePrice    float64
	Bid            float64
	Ask            float64
	ExpirationDate time.Time
	Position       PositionType
}

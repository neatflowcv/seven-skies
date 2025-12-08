package flow

import "time"

type Weather struct {
	Source       string
	TargetDate   time.Time
	ForecastDate time.Time
	Condition    string
	Temperature  float64 // Celsius
}

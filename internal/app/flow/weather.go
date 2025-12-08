package flow

import "time"

type Weather struct {
	Source       string
	TargetDate   time.Time
	ForecastDate time.Time
	Condition    string
	Temperature  float64 // Celsius
}

type DailyWeather struct {
	Date      time.Time
	High      float64 // Celsius
	Low       float64 // Celsius
	Condition string
}

package domain

import "time"

type Celsius float64

type Temperature struct {
	Value Celsius `json:"value"`
}

type WeatherSource string

const (
	WeatherSourceOpenWeather WeatherSource = "openweather"
	WeatherSourceKMA         WeatherSource = "kma"
)

type WeatherEvent struct {
	// Source is the source of the weather event. openweather, kma etc
	Source       WeatherSource    `json:"source"`
	TargetDate   time.Time        `json:"target_date"`
	ForecastDate time.Time        `json:"forecast_date"`
	Condition    WeatherCondition `json:"condition"`
	Temperature  Temperature      `json:"temperature,omitzero"`
}

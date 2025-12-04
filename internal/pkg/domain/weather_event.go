package domain

import "time"

type Celsius float64

type Temperature struct {
	Value Celsius
}

type WeatherEvent struct {
	Date        time.Time
	Condition   WeatherCondition `json:"condition"`
	Temperature *Temperature     `json:"temperature"`
	High        *Temperature     `json:"high,omitempty"`
	Low         *Temperature     `json:"low,omitempty"`
}

package domain

import "time"

type Celsius float64

type Temperature struct {
	Value Celsius `json:"value"`
}

type WeatherEvent struct {
	Date        time.Time        `json:"date"`
	Condition   WeatherCondition `json:"condition"`
	Temperature *Temperature     `json:"temperature,omitempty"`
	High        *Temperature     `json:"high,omitempty"`
	Low         *Temperature     `json:"low,omitempty"`
}

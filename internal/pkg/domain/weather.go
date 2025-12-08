package domain

import "time"

type Weather struct {
	id           string
	source       WeatherSource
	targetDate   time.Time
	forecastDate time.Time
	condition    WeatherCondition
	temperature  Temperature
}

func NewWeather(
	id string,
	source WeatherSource,
	targetDate time.Time,
	forecastDate time.Time,
	condition WeatherCondition,
	temperature Temperature,
) *Weather {
	return &Weather{
		id:           id,
		source:       source,
		targetDate:   targetDate,
		forecastDate: forecastDate,
		condition:    condition,
		temperature:  temperature,
	}
}

func (w *Weather) ID() string {
	return w.id
}

func (w *Weather) Source() WeatherSource {
	return w.source
}

func (w *Weather) TargetDate() time.Time {
	return w.targetDate
}

func (w *Weather) ForecastDate() time.Time {
	return w.forecastDate
}

func (w *Weather) Condition() WeatherCondition {
	return w.condition
}

func (w *Weather) Temperature() Temperature {
	return w.temperature
}

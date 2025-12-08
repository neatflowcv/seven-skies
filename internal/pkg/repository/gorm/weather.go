package gorm

import (
	"time"

	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
)

type Weather struct {
	ID           string `gorm:"primaryKey"`
	Source       string
	TargetDate   time.Time
	ForecastDate time.Time
	Condition    string
	Temperature  float64
}

func newModelWeather(weather *domain.Weather) *Weather {
	return &Weather{
		ID:           weather.ID(),
		Source:       string(weather.Source()),
		TargetDate:   weather.TargetDate(),
		ForecastDate: weather.ForecastDate(),
		Condition:    string(weather.Condition()),
		Temperature:  float64(weather.Temperature().Value),
	}
}

func (w *Weather) toDomain() *domain.Weather {
	return domain.NewWeather(
		w.ID,
		domain.WeatherSource(w.Source),
		w.TargetDate,
		w.ForecastDate,
		domain.WeatherCondition(w.Condition),
		domain.Temperature{Value: domain.Celsius(w.Temperature)},
	)
}

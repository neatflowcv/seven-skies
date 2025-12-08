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

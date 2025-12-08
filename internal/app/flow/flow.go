package flow

import (
	"context"
	"errors"
	"fmt"

	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
	"github.com/neatflowcv/seven-skies/internal/pkg/repository"
	"github.com/oklog/ulid/v2"
)

type Flow struct {
	repo repository.Repository
}

func NewFlow(repo repository.Repository) *Flow {
	return &Flow{
		repo: repo,
	}
}

func (f *Flow) CreateWeather(ctx context.Context, weather *Weather) error {
	id := ulid.Make().String()

	domainWeather := domain.NewWeather(
		id,
		domain.WeatherSource(weather.Source),
		weather.TargetDate,
		weather.ForecastDate,
		domain.WeatherCondition(weather.Condition),
		domain.Temperature{
			Value: domain.Celsius(weather.Temperature),
		},
	)

	err := f.repo.CreateWeather(ctx, domainWeather)
	if err != nil {
		if errors.Is(err, repository.ErrWeatherAlreadyExists) {
			return ErrWeatherAlreadyExists
		}

		return fmt.Errorf("error in create weather: %w", err)
	}

	return nil
}

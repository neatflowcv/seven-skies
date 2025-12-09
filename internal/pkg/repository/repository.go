package repository

import (
	"context"
	"time"

	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
)

//go:generate mockgen -typed -package=mocks -destination=mocks/repository.go . Repository

type Repository interface {
	CreateWeather(ctx context.Context, weather *domain.Weather) error
	ListWeathers(ctx context.Context, from, to time.Time) ([]*domain.Weather, error)
}

package repository

import (
	"context"
	"time"

	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
)

type Repository interface {
	CreateWeather(ctx context.Context, weather *domain.Weather) error
	ListWeathers(ctx context.Context, from, to time.Time) ([]*domain.Weather, error)
}

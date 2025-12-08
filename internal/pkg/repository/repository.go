package repository

import (
	"context"

	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
)

type Repository interface {
	CreateWeather(ctx context.Context, weather *domain.Weather) error
}

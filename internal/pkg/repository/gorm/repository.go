package gorm

import (
	"context"
	"fmt"

	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
	"github.com/neatflowcv/seven-skies/internal/pkg/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ repository.Repository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) //nolint:exhaustruct
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.AutoMigrate(&Weather{}) //nolint:exhaustruct
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) CreateWeather(ctx context.Context, weather *domain.Weather) error {
	modelWeather := newModelWeather(weather)

	err := gorm.G[Weather](r.db).Create(ctx, modelWeather)
	if err != nil {
		return fmt.Errorf("failed to create weather: %w", err)
	}

	return nil
}

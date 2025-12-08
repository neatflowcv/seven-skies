package flow

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

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

func (f *Flow) ListDailyWeathers(
	ctx context.Context,
	timezone string,
	now time.Time,
) ([]*DailyWeather, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("error in load location: %w", err)
	}

	const aWeek = 7

	localNow := now.In(location)

	from := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, location)
	to := from.AddDate(0, 0, aWeek).Add(-time.Nanosecond)

	weathers, err := f.repo.ListWeathers(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("error in list daily weathers: %w", err)
	}

	weathersByDate := f.groupByDate(weathers, location)

	var dailyWeathers []*DailyWeather

	for dateKey, weathers := range weathersByDate {
		newVar := f.mergeWeathers(dateKey, location, weathers)
		dailyWeathers = append(dailyWeathers, newVar)
	}

	sort.Slice(dailyWeathers, func(i, j int) bool {
		return dailyWeathers[i].Date.Before(dailyWeathers[j].Date)
	})

	return dailyWeathers, nil
}

func (*Flow) mergeWeathers(dateKey string, location *time.Location, weathers []*domain.Weather) *DailyWeather {
	date, _ := time.Parse("2006-01-02", dateKey)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)

	high := float64(weathers[0].Temperature().Value)
	low := float64(weathers[0].Temperature().Value)
	worstCondition := weathers[0].Condition()

	for _, w := range weathers {
		temp := float64(w.Temperature().Value)
		if temp > high {
			high = temp
		}

		if temp < low {
			low = temp
		}

		if w.Condition().IsWorseThan(worstCondition) {
			worstCondition = w.Condition()
		}
	}

	newVar := &DailyWeather{
		Date:      date,
		High:      high,
		Low:       low,
		Condition: string(worstCondition),
	}

	return newVar
}

func (*Flow) groupByDate(weathers []*domain.Weather, location *time.Location) map[string][]*domain.Weather {
	weathersByDate := make(map[string][]*domain.Weather)

	for _, weather := range weathers {
		localDate := weather.TargetDate().In(location)
		dateKey := time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 0, 0, 0, 0, location).Format("2006-01-02")
		weathersByDate[dateKey] = append(weathersByDate[dateKey], weather)
	}

	return weathersByDate
}

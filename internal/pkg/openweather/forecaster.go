package openweather

import (
	"context"
	"errors"
	"fmt"
	"log"

	"resty.dev/v3"
)

var (
	ErrForecast = errors.New("error in get forecast")
)

func Forecast(ctx context.Context, key string, lat float64, lon float64) (*ForecastResponse, error) {
	client := resty.New()

	defer func() {
		err := client.Close()
		if err != nil {
			log.Println("error in close client", err)
		}
	}()

	resp, err := client.R().
		SetContext(ctx).
		SetQueryParam("lat", fmt.Sprintf("%f", lat)).
		SetQueryParam("lon", fmt.Sprintf("%f", lon)).
		SetQueryParam("appid", key).
		SetQueryParam("units", "metric").
		SetResult(&ForecastResponse{}). //nolint:exhaustruct
		Get("https://api.openweathermap.org/data/2.5/forecast")
	if err != nil {
		return nil, fmt.Errorf("error in get forecast: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("%w: %s, %s", ErrForecast, resp.Status(), resp.String())
	}

	return resp.Result().(*ForecastResponse), nil //nolint:forcetypeassert
}

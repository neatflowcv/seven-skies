package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/neatflowcv/seven-skies/internal/pkg/broker"
	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
	"github.com/neatflowcv/seven-skies/internal/pkg/openweather"
)

type Handler struct {
	broker broker.Broker
	key    string
	lat    float64
	lon    float64
}

func NewHandler(broker broker.Broker, key string, lat float64, lon float64) *Handler {
	return &Handler{
		broker: broker,
		key:    key,
		lat:    lat,
		lon:    lon,
	}
}

func (h *Handler) Handle(ctx context.Context) error {
	log.Println("get forecast")

	forecast, err := openweather.Forecast(ctx, h.key, h.lat, h.lon)
	if err != nil {
		log.Println("error in get forecast", err)

		return err
	}

	log.Println("get forecast successfully", forecast)

	for _, list := range forecast.List {
		weatherEvent := domain.WeatherEvent{
			Date:        time.Unix(int64(list.Dt), 0),
			Condition:   decideCondition(list.Weather[0].Main),
			Temperature: &domain.Temperature{Value: domain.Celsius(list.Main.Temp)},
			High:        &domain.Temperature{Value: domain.Celsius(list.Main.TempMax)},
			Low:         &domain.Temperature{Value: domain.Celsius(list.Main.TempMin)},
		}

		message, err := json.Marshal(weatherEvent)
		if err != nil {
			log.Println("error in marshal weather event", err)

			continue
		}

		err = h.broker.Publish(ctx, "SEVEN_SKIES_SUBJECT.OPENWEATHER", message)
		if err != nil {
			log.Println("error in publish weather event", err)

			continue
		}
	}

	return nil
}

func decideCondition(main string) domain.WeatherCondition {
	switch main {
	case "Thunderstorm", "Drizzle", "Rain":
		return domain.WeatherConditionRain

	case "Snow":
		return domain.WeatherConditionSnow

	case "Mist", "Smoke", "Haze", "Dust", "Fog", "Sand", "Ash", "Squall", "Tornado", "Clouds":
		return domain.WeatherConditionCloudy

	case "Clear":
		return domain.WeatherConditionClear

	default:
		return domain.WeatherConditionUnknown
	}
}

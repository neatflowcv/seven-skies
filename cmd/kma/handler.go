package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/neatflowcv/seven-skies/internal/pkg/broker"
	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
	"github.com/neatflowcv/seven-skies/internal/pkg/kma"
)

type Handler struct {
	broker broker.Broker
	key    string
	nx     int
	ny     int
}

func NewHandler(b broker.Broker, key string, nx int, ny int) *Handler {
	return &Handler{
		broker: b,
		key:    key,
		nx:     nx,
		ny:     ny,
	}
}

type TemporalData struct {
	Temperature   string
	Sky           string
	Precipitation string
	ForecastDate  string
	TargetDate    string
}

func buildTemporalData(resp *kma.ForecastResponse) []TemporalData {
	dates := map[string]*TemporalData{}
	for _, item := range resp.Response.Body.Items.Item {
		data, ok := dates[item.FcstDate+item.FcstTime]
		if !ok {
			data = &TemporalData{
				ForecastDate:  item.BaseDate + item.BaseTime,
				TargetDate:    item.FcstDate + item.FcstTime,
				Temperature:   "",
				Sky:           "",
				Precipitation: "",
			}
		}

		switch item.Category {
		case "T1H", "TMP": // 온도
			data.Temperature = item.FcstValue
		case "SKY": // 하늘상태
			data.Sky = item.FcstValue
		case "PTY": // 강수형태
			data.Precipitation = item.FcstValue
		}

		dates[item.FcstDate+item.FcstTime] = data
	}

	var ret []TemporalData
	for _, data := range dates {
		ret = append(ret, *data)
	}

	return ret
}

func parseAndValidateTemporalData(data TemporalData) (*domain.WeatherEvent, bool) {
	targetTime, err := parseKMAForecastDate(data.TargetDate)
	if err != nil {
		log.Println("error in parse forecast datetime", err)

		return nil, false
	}

	forecastTime, err := parseKMAForecastDate(data.ForecastDate)
	if err != nil {
		log.Println("error in parse forecast datetime", err)

		return nil, false
	}

	temperature, err := strconv.ParseInt(data.Temperature, 10, 64)
	if err != nil {
		log.Println("error in parse temperature", err)

		return nil, false
	}

	condition := decideCondition(&data)

	// 값 검증: 비어있는 값이 있으면 건너뛰기
	if forecastTime.IsZero() {
		log.Printf("skip weather event: forecastTime is empty (targetDate: %v)", targetTime)

		return nil, false
	}

	if targetTime.IsZero() {
		log.Printf("skip weather event: targetTime is empty (forecastDate: %v)", forecastTime)

		return nil, false
	}

	if condition == "" {
		log.Printf("skip weather event: condition is empty (forecastDate: %v, targetDate: %v)", forecastTime, targetTime)

		return nil, false
	}

	if temperature == 0 {
		log.Printf("skip weather event: temperature is empty (forecastDate: %v, targetDate: %v)", forecastTime, targetTime)

		return nil, false
	}

	return &domain.WeatherEvent{
		Source:       domain.WeatherSourceKMA,
		ForecastDate: forecastTime,
		TargetDate:   targetTime,
		Condition:    condition,
		Temperature: domain.Temperature{
			Value: domain.Celsius(temperature),
		},
	}, true
}

func convertTemporalDataToWeatherEvents(datas []TemporalData) []*domain.WeatherEvent {
	var weatherEvents []*domain.WeatherEvent

	for _, data := range datas {
		weatherEvent, ok := parseAndValidateTemporalData(data)
		if !ok {
			continue
		}

		weatherEvents = append(weatherEvents, weatherEvent)
	}

	return weatherEvents
}

// HandleUltraShort는 초단기 예보를 조회해 이벤트로 발행한다.
func (h *Handler) HandleUltraShort(ctx context.Context) error {
	log.Println("get ultra short forecast")

	now := time.Now()

	resp, err := kma.GetUltraShortForecast(ctx, h.key, now, h.nx, h.ny)
	if err != nil {
		return fmt.Errorf("error in get ultra short forecast: %w", err)
	}

	datas := buildTemporalData(resp)
	weatherEvents := convertTemporalDataToWeatherEvents(datas)

	for _, weatherEvent := range weatherEvents {
		message, err := json.Marshal(weatherEvent)
		if err != nil {
			log.Println("error in marshal weather event", err)

			continue
		}

		err = h.broker.Publish(ctx, "SEVEN_SKIES_SUBJECT.KMA", message)
		if err != nil {
			log.Println("error in publish weather event", err)

			continue
		}
	}

	return nil
}

func decideCondition(data *TemporalData) domain.WeatherCondition {
	if data.Precipitation == "0" && data.Sky == "1" {
		return domain.WeatherConditionClear
	}

	switch data.Precipitation {
	case "1", "5":
		return domain.WeatherConditionRain
	case "2", "6":
		return domain.WeatherConditionRainSnow
	case "3", "7":
		return domain.WeatherConditionSnow
	}

	switch data.Sky {
	case "3", "4":
		return domain.WeatherConditionCloudy
	}

	return domain.WeatherConditionUnknown
}

// HandleShortTerm는 단기 예보를 조회해 이벤트로 발행한다.
func (h *Handler) HandleShortTerm(ctx context.Context) error {
	log.Println("get short term forecast")

	now := time.Now()

	resp, err := kma.GetShortTermForecast(ctx, h.key, now, h.nx, h.ny)
	if err != nil {
		return fmt.Errorf("error in get short term forecast: %w", err)
	}

	datas := buildTemporalData(resp)
	weatherEvents := convertTemporalDataToWeatherEvents(datas)

	for _, weatherEvent := range weatherEvents {
		message, err := json.Marshal(weatherEvent)
		if err != nil {
			log.Println("error in marshal weather event", err)

			continue
		}

		err = h.broker.Publish(ctx, "SEVEN_SKIES_SUBJECT.KMA", message)
		if err != nil {
			log.Println("error in publish weather event", err)

			continue
		}
	}

	return nil
}

func parseKMAForecastDate(dateStr string) (time.Time, error) {
	layout := "200601021504"

	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}

	t, err := time.ParseInLocation(layout, dateStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("error in parse datetime: %w", err)
	}

	return t, nil
}

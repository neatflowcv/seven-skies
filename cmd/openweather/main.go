package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"github.com/neatflowcv/seven-skies/internal/pkg/broker"
	"github.com/neatflowcv/seven-skies/internal/pkg/broker/nats"
	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
	"github.com/neatflowcv/seven-skies/internal/pkg/openweather"
)

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return info.Main.Version
}

func main() {
	log.Println("version", version())

	godotenv.Load()

	key := os.Getenv("OPENWEATHER_API_KEY")
	if key == "" {
		log.Panic("OPENWEATHER_API_KEY is not set")
	}

	lat := os.Getenv("OPENWEATHER_LAT")
	if lat == "" {
		log.Panic("OPENWEATHER_LAT is not set")
	}

	lon := os.Getenv("OPENWEATHER_LON")
	if lon == "" {
		log.Panic("OPENWEATHER_LON is not set")
	}

	NATS_URL := os.Getenv("NATS_URL")
	if NATS_URL == "" {
		log.Panic("NATS_URL is not set")
	}

	broker, err := nats.NewBroker(context.Background(), NATS_URL, "SEVEN_SKIES_STREAM", []string{"SEVEN_SKIES_SUBJECT.>"})
	if err != nil {
		log.Panic("error in create broker", err)
	}
	defer broker.Close()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Panic("error in create scheduler", err)
	}

	defer func() {
		err := scheduler.Shutdown()
		if err != nil {
			log.Println("error in shutdown scheduler", err)
		}
	}()

	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Panic("error in parse lat", err)
	}

	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		log.Panic("error in parse lon", err)
	}

	err = Task(context.Background(), broker, key, latFloat, lonFloat)
	if err != nil {
		log.Panic("error in task", err)
	}

	job, err := scheduler.NewJob(
		gocron.CronJob("5 */2 * * *", false), // openweather에서 2시간마다 데이터가 업데이트 된다.
		// 5분을 준 이유는 업데이트가 바로 된다는 보장이 없어서
		gocron.NewTask(
			Task,
			context.Background(),
			broker,
			key,
			latFloat,
			lonFloat,
		),
	)
	if err != nil {
		log.Panic("error in create job", err)
	}

	scheduler.Start()

	nextRun, err := job.NextRun()
	if err != nil {
		log.Panic("error in get next run", err)
	}

	log.Println("next run", nextRun)

	select {}
}

func Task(ctx context.Context, broker broker.Broker, key string, lat float64, lon float64) error {
	log.Println("get forecast")

	forecast, err := openweather.Forecast(context.Background(), key, lat, lon)
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

		err = broker.Publish(ctx, "SEVEN_SKIES_SUBJECT.OPENWEATHER", message)
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

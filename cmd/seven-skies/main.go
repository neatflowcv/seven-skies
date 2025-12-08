package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
	"github.com/neatflowcv/seven-skies/internal/app/flow"
	"github.com/neatflowcv/seven-skies/internal/pkg/broker/nats"
	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
	"github.com/neatflowcv/seven-skies/internal/pkg/repository/gorm"
)

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return info.Main.Version
}

type Config struct {
	NATS_URL     string
	DATABASE_DSN string
}

func NewConfig() *Config {
	NATS_URL := os.Getenv("NATS_URL")
	if NATS_URL == "" {
		log.Panic("NATS_URL is not set")
	}

	DATABASE_DSN := os.Getenv("DATABASE_DSN")
	if DATABASE_DSN == "" {
		log.Panic("DATABASE_DSN is not set")
	}

	return &Config{
		NATS_URL:     NATS_URL,
		DATABASE_DSN: DATABASE_DSN,
	}
}

func main() {
	log.Println("version", version())

	err := godotenv.Load()
	if err == nil {
		log.Println("env loaded")
	}

	config := NewConfig()

	ctx := context.Background()

	broker, err := nats.NewBroker(ctx, config.NATS_URL, "SEVEN_SKIES_STREAM", []string{"SEVEN_SKIES_SUBJECT.>"})
	if err != nil {
		log.Panic("error in create broker", err)
	}
	defer broker.Close()

	repo, err := gorm.NewRepository(config.DATABASE_DSN)
	if err != nil {
		log.Panic("error in create repository", err)
	}

	service := flow.NewFlow(repo)

	err = broker.Subscribe(ctx, "SEVEN_SKIES_DURABLE", func(topic string, message []byte) error {
		var event domain.WeatherEvent

		err := json.Unmarshal(message, &event)
		if err != nil {
			log.Println("error in unmarshal", err)

			return fmt.Errorf("error in unmarshal: %w", err)
		}

		err = service.CreateWeather(ctx, &flow.Weather{
			Source:       string(event.Source),
			TargetDate:   event.TargetDate,
			ForecastDate: event.ForecastDate,
			Condition:    string(event.Condition),
			Temperature:  float64(event.Temperature.Value),
		})
		if err != nil {
			return fmt.Errorf("error in create weather: %w", err)
		}

		return nil
	})
	if err != nil {
		log.Panic("error in subscribe", err)
	}

	select {}
}

package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"github.com/neatflowcv/seven-skies/internal/pkg/broker/nats"
)

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return info.Main.Version
}

type Config struct {
	Key     string
	Lat     float64
	Lon     float64
	NATSURL string
}

func NewConfig() *Config {
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

	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Panic("error in parse lat", err)
	}

	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		log.Panic("error in parse lon", err)
	}

	return &Config{
		Key:     key,
		Lat:     latFloat,
		Lon:     lonFloat,
		NATSURL: NATS_URL,
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

	broker, err := nats.NewBroker(ctx, config.NATSURL, "SEVEN_SKIES_STREAM", []string{"SEVEN_SKIES_SUBJECT.>"})
	if err != nil {
		log.Panic("error in create broker", err)
	}
	defer broker.Close()

	handler := NewHandler(broker, config.Key, config.Lat, config.Lon)

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

	err = handler.Handle(ctx)
	if err != nil {
		log.Panic("error in task", err)
	}

	job, err := scheduler.NewJob(
		gocron.CronJob("5 */2 * * *", false), // openweather에서 2시간마다 데이터가 업데이트 된다.
		// 5분을 준 이유는 업데이트가 바로 된다는 보장이 없어서
		gocron.NewTask(
			handler.Handle,
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

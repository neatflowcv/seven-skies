package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
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

	log.Println("version", version())

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

	err = Task(key, latFloat, lonFloat)
	if err != nil {
		log.Panic("error in task", err)
	}

	job, err := scheduler.NewJob(
		gocron.CronJob("5 */2 * * *", false), // openweather에서 2시간마다 데이터가 업데이트 된다.
		// 5분을 준 이유는 업데이트가 바로 된다는 보장이 없어서
		gocron.NewTask(
			Task,
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

func Task(key string, lat float64, lon float64) error {
	log.Println("get forecast")
	forecast, err := openweather.Forecast(context.Background(), key, lat, lon)
	if err != nil {
		log.Println("error in get forecast", err)
		return err
	}

	log.Println("get forecast successfully", forecast)

	return nil
}

package main

import (
	"context"
	"fmt"
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
	NX      int
	NY      int
	NATSURL string
}

func NewConfig() *Config {
	key := os.Getenv("KMA_API_KEY")
	if key == "" {
		log.Panic("KMA_API_KEY is not set")
	}

	nx := os.Getenv("KMA_NX")
	if nx == "" {
		log.Panic("KMA_NX is not set")
	}

	ny := os.Getenv("KMA_NY")
	if ny == "" {
		log.Panic("KMA_NY is not set")
	}

	NATS_URL := os.Getenv("NATS_URL")
	if NATS_URL == "" {
		log.Panic("NATS_URL is not set")
	}

	nxInt, err := strconv.Atoi(nx)
	if err != nil {
		log.Panic("error in parse nx", err)
	}

	nyInt, err := strconv.Atoi(ny)
	if err != nil {
		log.Panic("error in parse ny", err)
	}

	return &Config{
		Key:     key,
		NX:      nxInt,
		NY:      nyInt,
		NATSURL: NATS_URL,
	}
}

func setupBroker(ctx context.Context, natsURL string) (*nats.Broker, error) {
	broker, err := nats.NewBroker(ctx, natsURL, "SEVEN_SKIES_STREAM", []string{"SEVEN_SKIES_SUBJECT.>"})
	if err != nil {
		return nil, fmt.Errorf("error in create broker: %w", err)
	}

	return broker, nil
}

//nolint:ireturn // gocron.Scheduler is an interface by design
func setupScheduler() (gocron.Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("error in create scheduler: %w", err)
	}

	return scheduler, nil
}

func runInitialTasks(ctx context.Context, handler *Handler) error {
	err := handler.HandleUltraShort(ctx)
	if err != nil {
		return err
	}

	err = handler.HandleShortTerm(ctx)
	if err != nil {
		return err
	}

	return nil
}

//nolint:ireturn // gocron.Job is an interface by design
func setupJobs(scheduler gocron.Scheduler, handler *Handler) (gocron.Job, gocron.Job, error) {
	ultraJob, err := scheduler.NewJob(
		gocron.CronJob("50 * * * *", false),
		gocron.NewTask(handler.HandleUltraShort),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error in create ultra short job: %w", err)
	}

	shortJob, err := scheduler.NewJob(
		gocron.CronJob("15 2,5,8,11,14,17,20,23 * * *", false),
		gocron.NewTask(handler.HandleShortTerm),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error in create short term job: %w", err)
	}

	return ultraJob, shortJob, nil
}

func logNextRuns(ultraJob, shortJob gocron.Job) error {
	ultraNext, err := ultraJob.NextRun()
	if err != nil {
		return fmt.Errorf("error in get ultra short next run: %w", err)
	}

	shortNext, err := shortJob.NextRun()
	if err != nil {
		return fmt.Errorf("error in get short term next run: %w", err)
	}

	log.Println("ultra short next run", ultraNext)
	log.Println("short term next run", shortNext)

	return nil
}

func main() {
	log.Println("version", version())

	err := godotenv.Load()
	if err == nil {
		log.Println("env loaded")
	}

	config := NewConfig()
	ctx := context.Background()

	broker, err := setupBroker(ctx, config.NATSURL)
	if err != nil {
		log.Panic("error in create broker", err)
	}
	defer broker.Close()

	handler := NewHandler(broker, config.Key, config.NX, config.NY)

	scheduler, err := setupScheduler()
	if err != nil {
		log.Panic("error in create scheduler", err)
	}

	defer func() {
		err := scheduler.Shutdown()
		if err != nil {
			log.Println("error in shutdown scheduler", err)
		}
	}()

	err = runInitialTasks(ctx, handler)
	if err != nil {
		log.Panic("error in initial tasks", err)
	}

	ultraJob, shortJob, err := setupJobs(scheduler, handler)
	if err != nil {
		log.Panic("error in create jobs", err)
	}

	scheduler.Start()

	err = logNextRuns(ultraJob, shortJob)
	if err != nil {
		log.Panic("error in get next runs", err)
	}

	select {}
}

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	eventservice "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/service"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	migrate(config.Storage.Type, config.Storage.SQLDriver, config.Storage.DSN)

	log, err := logger.New(config.Logger.Level, config.Logger.OutputPaths, config.Logger.ErrorOutputPaths)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	eventService, err := eventservice.New(ctx, config.Storage.Type, config.Storage.SQLDriver, config.Storage.DSN)
	if err != nil {
		log.Fatal(err)
	}

	calendar := app.New(log, eventService)
	server := internalhttp.NewServer(log, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

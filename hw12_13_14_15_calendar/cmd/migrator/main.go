package main

import (
	"context"
	"flag"
	"log"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	l "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config.NewCalendarConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := l.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	strg := storage.New(ctx, logger, config.Storage)
	if strg.Connection != nil {
		migrate(strg.Connection.DB, logger)
	}
}

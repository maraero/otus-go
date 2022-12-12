package main

import (
	_ "github.com/lib/pq"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	_ "github.com/maraero/otus-go/hw12_13_14_15_calendar/migrations"
	"github.com/pressly/goose/v3"
)

func migrate(logger *logger.Log, c config.Storage) {
	if c.Type == config.StorageInMemory {
		return
	}

	db, err := goose.OpenDBWithDriver(c.Driver, c.DSN)
	if err != nil {
		logger.Fatal("goose: failed to open DB: %w\n", err)
	}
	logger.Info("Database connection was created\n")

	defer func() {
		if err := db.Close(); err != nil {
			logger.Fatal("goose: failed to close DB: %w\n", err)
		}
		logger.Info("Database connection was closed\n")
	}()

	logger.Info("Start migrating database\n")
	if err := goose.Up(db, "."); err != nil {
		logger.Error("goose Up error: %w", err)
	}
	logger.Info("Complete migrating database\n")
}

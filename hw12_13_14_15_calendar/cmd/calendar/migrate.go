package main

import (
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/pressly/goose/v3"
)

func migrate(log *logger.Log, c config.Storage) {
	if c.Type == config.StorageInMemory {
		return
	}

	db, err := goose.OpenDBWithDriver(c.Database, c.DSN)
	if err != nil {
		log.Fatal("goose: failed to open DB: %w\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("goose: failed to close DB: %w\n", err)
		}
	}()

	if err := goose.Up(db, "../../migrations"); err != nil {
		log.Error("goose Up error: %w", err)
		return
	}
}

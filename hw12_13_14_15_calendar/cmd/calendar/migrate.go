package main

import (
	"fmt"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/pressly/goose/v3"
)

func migrate(log *logger.Log, storageType string, sqlDriver string, dsn string) {
	if storageType == config.StorageInMemory {
		return
	}

	driver, err := getDBDriverBySQLDriver(sqlDriver)
	if err != nil {
		log.Fatal(err)
	}

	db, err := goose.OpenDBWithDriver(driver, dsn)
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

func getDBDriverBySQLDriver(sqlDriver string) (string, error) {
	if sqlDriver == "pgx" {
		return "postgres", nil
	}
	return "", fmt.Errorf("Unknown sql driver: %v", sqlDriver)
}
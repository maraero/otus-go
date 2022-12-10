package main

import (
	"fmt"
	"log"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/pressly/goose/v3"
)

func migrate(storageType string, sqlDriver string, DSN string) {
	if storageType == config.StorageInMemory {
		return
	}

	driver, err := getDBDriverBySqlDriver(sqlDriver)
	if err != nil {
		log.Fatal(err)
	}

	db, err := goose.OpenDBWithDriver(driver, DSN)
	if err != nil {
		log.Fatal("goose: failed to open DB: %w\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("goose: failed to close DB: %w\n", err)
		}
	}()

	if err := goose.Up(db, "../../migrations"); err != nil {
		log.Fatal("goose Up error: %w", err)
	}
}

func getDBDriverBySqlDriver(sqlDriver string) (string, error) {
	if sqlDriver == "pgx" {
		return "postgres", nil
	}
	return "", fmt.Errorf("Unknown sql driver: %v", sqlDriver)
}

package main

import (
	"database/sql"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	_ "github.com/maraero/otus-go/hw12_13_14_15_calendar/migrations"
	"github.com/pressly/goose/v3"
)

func migrate(db *sql.DB, logger *logger.Log) {
	logger.Info("Start migrating database\n")
	if err := goose.Up(db, "."); err != nil {
		logger.Error("goose Up error: %w", err)
	}
	logger.Info("Complete migrating database\n")
}

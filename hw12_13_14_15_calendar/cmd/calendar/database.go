package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
)

const (
	MaxOpenConns    = 25
	MaxIdleConns    = 25
	ConnMaxLifetime = time.Minute
)

func newDBConnection(ctx context.Context, c config.Storage) *sqlx.DB {
	if c.Type != config.StorageSQL {
		return nil
	}

	conn, err := connect(ctx, c.Driver, c.DSN)
	if err != nil {
		log.Fatal("DB connection error: %w", err)
	}

	return conn
}

func connect(ctx context.Context, driver string, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify connection to database: %w", err)
	}

	return db, nil
}

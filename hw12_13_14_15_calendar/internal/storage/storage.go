package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

func New(ctx context.Context, logger *logger.Log, c config.Storage) *Storage {
	if c.Type != config.StorageSQL {
		return &Storage{
			Source:     StorageInMemory,
			Connection: nil,
		}
	}

	conn, err := connectToDB(ctx, c.Driver, c.DSN)
	if err != nil {
		logger.Fatal("DB connection error: %w", err)
	}

	return &Storage{
		Source:     StorageSQL,
		Connection: conn,
	}
}

func connectToDB(ctx context.Context, driver string, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify connection to database: %w", err)
	}

	return db, nil
}

package storage

import "time"

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = time.Minute
)

type Source = string

const (
	StorageInMemory Source = "in-memory"
	StorageSQL      Source = "sql"
)

package storage

import "github.com/jmoiron/sqlx"

type StorageType = string

type Storage struct {
	Source     StorageType
	Connection *sqlx.DB
}

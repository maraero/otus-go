package storage

import "github.com/jmoiron/sqlx"

type Source = string

type Storage struct {
	Source     Source
	Connection *sqlx.DB
}

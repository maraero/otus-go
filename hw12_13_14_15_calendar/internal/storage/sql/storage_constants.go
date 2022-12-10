package sqlstorage

import "time"

const (
	MaxOpenConns    = 25
	MaxIdleConns    = 25
	ConnMaxLifetime = time.Minute
)

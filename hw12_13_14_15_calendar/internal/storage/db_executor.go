package storage

import (
	"context"
	"database/sql"
)

type DBExecutor interface {
	ExecContext(ctx context.Context, query string, arg ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
